package demo_import

import (
	"encoding/json"
	"fmt"
	"github.com/opengovern/og-util/pkg/integration"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	Parallelism = 20
	BatchLimit  = 10000
	ScrollTime  = "10m"
)

type Config struct {
	DemoDataURL       string
	OpenSSLPassword   string
	ElasticsearchUser string
	ElasticsearchPass string
	ElasticsearchAddr string
}

type Integration struct {
	IntegrationID   string
	ProviderID      string
	Name            string
	IntegrationType integration.Type
	Annotations     map[string]string
	Labels          map[string]string
}

func LoadDemoData(cfg Config, logger *zap.Logger) ([]Integration, error) {
	logger.Info("Starting data loading process with provided config...")

	wd, err := os.Getwd()
	if err != nil {
		// Log error before returning
		logger.Error("Failed to get current working directory", zap.Error(err))
		return nil, fmt.Errorf("failed to get current working directory: %w", err)
	}
	workDir := wd

	logger.Info("Using working directory", zap.String("path", workDir))

	encryptedFilePath := filepath.Join(workDir, "demo_data.tar.gz.enc")
	decryptedFilePath := filepath.Join(workDir, "demo_data.tar.gz")
	// Assuming tar extracts relative paths inside the archive into the current dir (workDir).
	inputPathForDump := filepath.Join(workDir, "/demo-data/es-demo/")
	integrationsJsonFilePath := filepath.Join(workDir, "/demo-data/integrations.json")

	// --- Ensure cleanup of intermediate files ---
	defer func() {
		logger.Info("Cleaning up intermediate file", zap.String("path", encryptedFilePath))
		err := os.Remove(encryptedFilePath)
		// Log only if it's an actual error other than "not exist"
		if err != nil && !os.IsNotExist(err) {
			logger.Warn("Failed to remove encrypted file", zap.String("path", encryptedFilePath), zap.Error(err))
		}
	}()
	defer func() {
		logger.Info("Cleaning up intermediate file", zap.String("path", decryptedFilePath))
		err := os.Remove(decryptedFilePath)
		if err != nil && !os.IsNotExist(err) {
			logger.Warn("Failed to remove decrypted file", zap.String("path", decryptedFilePath), zap.Error(err))
		}
	}()

	// --- 1. Download the file ---
	logger.Info("Downloading data",
		zap.String("url", cfg.DemoDataURL),
		zap.String("destination", encryptedFilePath))
	err = downloadFile(encryptedFilePath, cfg.DemoDataURL, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	logger.Info("Successfully downloaded file", zap.String("path", encryptedFilePath))

	// --- 2. Decrypt the file ---
	logger.Info("Decrypting file",
		zap.String("input", encryptedFilePath),
		zap.String("output", decryptedFilePath))
	opensslArgs := []string{
		"enc", "-d", "-aes-256-cbc", "-md", "md5",
		"-pass", "pass:" + cfg.OpenSSLPassword,
		"-base64",
		"-in", encryptedFilePath,
		"-out", decryptedFilePath,
	}
	err = runCommand(workDir, "openssl", logger, opensslArgs...) // Pass logger
	if err != nil {
		return nil, fmt.Errorf("failed to run openssl command: %w", err)
	}
	logger.Info("Successfully decrypted file", zap.String("path", decryptedFilePath))

	// --- 3. Extract the tarball ---
	logger.Info("Extracting tarball",
		zap.String("input", decryptedFilePath),
		zap.String("destinationDir", workDir))
	tarArgs := []string{"-xvf", decryptedFilePath}
	err = runCommand(workDir, "tar", logger, tarArgs...) // Pass logger
	if err != nil {
		return nil, fmt.Errorf("failed to run tar command: %w", err)
	}
	logger.Info("Successfully extracted tarball", zap.String("path", decryptedFilePath))

	// load integration file
	integrations, err := loadIntegrationsFromJSON(integrationsJsonFilePath)
	if err != nil {
		logger.Error("Failed to load integrations.", zap.Error(err))
		return nil, fmt.Errorf("failed to load integrations.: %w", err)
	}

	logger.Info("Successfully loaded integrations from json.", zap.Int("integrations", len(integrations)))

	// --- 4. Construct the new Elasticsearch address ---
	cleanAddress := strings.TrimPrefix(cfg.ElasticsearchAddr, "https://")
	newElasticsearchAddress := fmt.Sprintf("https://%s:%s@%s", cfg.ElasticsearchUser, cfg.ElasticsearchPass, cleanAddress)
	// Use a separate field for the address without credentials for safer logging
	logger.Info("Constructed Elasticsearch target address", zap.String("address", "https://****:****@"+cleanAddress))

	// --- 5. Run multielasticdump ---
	logger.Info("Preparing to run multielasticdump",
		zap.String("inputDir", inputPathForDump),
		zap.String("outputAddr", "https://****:****@"+cleanAddress), // Log sanitized address
		zap.Int("parallelism", Parallelism),
		zap.Int("limit", BatchLimit),
		zap.String("scrollTime", ScrollTime))

	dumpArgs := []string{
		"--direction=load",
		"--input=" + inputPathForDump,
		"--output=" + newElasticsearchAddress, // Use the full address here
		fmt.Sprintf("--parallel=%d", Parallelism),
		fmt.Sprintf("--limit=%d", BatchLimit),
		"--scrollTime=" + ScrollTime,
	}
	dumpArgs = append(dumpArgs, "--ignoreTemplate=true")

	cmd := exec.Command("multielasticdump", dumpArgs...)
	cmd.Stdout = os.Stdout // Keep piping command output for visibility
	cmd.Stderr = os.Stderr
	cmd.Dir = workDir // Run the command from the working directory

	// Set environment variable specifically for this command
	cmd.Env = append(os.Environ(), "NODE_TLS_REJECT_UNAUTHORIZED=0")

	logger.Info("Executing multielasticdump command",
		zap.String("command", "multielasticdump"),
		zap.Strings("args", dumpArgs), // Log arguments separately
		zap.String("workDir", workDir),
		zap.Strings("env_override", []string{"NODE_TLS_REJECT_UNAUTHORIZED=0"}))

	err = cmd.Run()
	if err != nil {
		logger.Error("multielasticdump command failed", zap.Error(err))
		return nil, fmt.Errorf("failed to run multielasticdump command: %w", err)
	}
	logger.Info("Successfully ran multielasticdump.")

	logger.Info("Data loading process completed successfully.")
	return integrations, nil
}

// downloadFile downloads a file from a URL and saves it locally, using zap for logging.
func downloadFile(filepath string, url string, logger *zap.Logger) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		logger.Error("Failed to create file for download", zap.String("path", filepath), zap.Error(err))
		return fmt.Errorf("failed to create file %s: %w", filepath, err)
	}
	defer out.Close() // Ensure file is closed

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		_ = os.Remove(filepath) // Clean up empty file
		logger.Error("HTTP GET request failed", zap.String("url", url), zap.Error(err))
		return fmt.Errorf("failed to get URL %s: %w", url, err)
	}
	defer resp.Body.Close() // Ensure response body is closed

	// Check server response
	if resp.StatusCode != http.StatusOK {
		_ = os.Remove(filepath) // Clean up empty file
		bodyBytes, _ := io.ReadAll(resp.Body)
		logger.Error("Received non-OK HTTP status",
			zap.Int("statusCode", resp.StatusCode),
			zap.String("url", url),
			zap.ByteString("responseBody", bodyBytes)) // Log response body if possible
		return fmt.Errorf("bad status code %d from %s", resp.StatusCode, url)
	}

	// Write the body to file
	bytesWritten, err := io.Copy(out, resp.Body)
	if err != nil {
		_ = os.Remove(filepath) // Clean up potentially partially written file
		logger.Error("Failed to write response body to file",
			zap.String("path", filepath),
			zap.Error(err))
		return fmt.Errorf("failed to write response body to file %s: %w", filepath, err)
	}
	logger.Debug("Successfully wrote file", zap.Int64("bytesWritten", bytesWritten), zap.String("path", filepath))

	return nil
}

// runCommand executes an external command and waits for it to finish, using zap for logging.
// It runs the command in the specified working directory.
// It redirects the command's stdout and stderr to the Go program's stderr.
func runCommand(workDir string, name string, logger *zap.Logger, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = workDir      // Set the working directory for the command
	cmd.Stdout = os.Stderr // Keep output visible
	cmd.Stderr = os.Stderr

	logger.Info("Executing command",
		zap.String("command", name),
		zap.Strings("args", args),
		zap.String("workDir", workDir))

	err := cmd.Run()
	if err != nil {
		// Log the error before returning it
		logger.Error("Command execution failed",
			zap.String("command", name),
			zap.Strings("args", args),
			zap.String("workDir", workDir),
			zap.Error(err)) // Include the underlying error from exec.Run
		return fmt.Errorf("command '%s %s' failed in %s: %w", name, strings.Join(args, " "), workDir, err)
	}
	logger.Debug("Command executed successfully", zap.String("command", name))
	return nil
}

func loadIntegrationsFromJSON(filePath string) ([]Integration, error) {
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	if len(jsonData) == 0 {
		return []Integration{}, nil
	}

	var integrations []Integration

	err = json.Unmarshal(jsonData, &integrations)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON from %s: %w", filePath, err)
	}

	return integrations, nil
}
