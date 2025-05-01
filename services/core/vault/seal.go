// vault/seal.go
package vault

import (
	"context"
	"errors" // Import errors
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"github.com/opengovern/og-util/pkg/vault"
	"github.com/opengovern/opensecurity/services/core/config"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	componentName = "VaultSealHandler"

	unsealCheckInterval  = 1 * time.Minute
	initialUnsealTimeout = 5 * time.Minute
)

// Specific errors returned by the handler
var (
	ErrVaultInitCheckFailed    = errors.New("failed checking vault init status")
	ErrVaultSecretUpdateFailed = errors.New("failed to update existing vault keys secret after init")
	ErrVaultSecretCreateFailed = errors.New("failed to create vault keys secret after init")
	ErrVaultSecretGetFailed    = errors.New("failed to get vault keys secret for verification")
	ErrVaultKeysMissing        = errors.New("vault keys secret is missing") // Inconsistent state
	ErrUnsealKeysGetFailed     = errors.New("unseal checker failed to get vault keys secret")
	ErrUnsealNoKeysFound       = errors.New("unseal checker found no keys in secret")
	ErrUnsealTimeout           = errors.New("timeout waiting for initial vault unseal signal")
	ErrUnsealCheckerExited     = errors.New("unseal checker exited prematurely before signaling success")
)

// REMOVED package-level var secretName
// var secretName = "vault-unseal-keys"

type SealHandler struct {
	logger           *zap.Logger
	cfg              config.Config
	vaultSealHandler *vault.HashiCorpVaultSealHandler
	kubeClientset    kubernetes.Interface
	secretName       string // <<< Field definition is correct
}

// NewSealHandler: **** MODIFY SIGNATURE and BODY ****
func NewSealHandler(ctx context.Context, logger *zap.Logger, cfg config.Config, secretNameIn string) (*SealHandler, error) { // <<< ADD secretNameIn parameter
	componentLogger := logger.With(
		zap.String("component", componentName),
		zap.String("namespace", cfg.OpengovernanceNamespace),
	)
	componentLogger.Debug("Initializing VaultSealHandler...")

	if cfg.OpengovernanceNamespace == "" {
		return nil, fmt.Errorf("opengovernance namespace is required in configuration")
	}

	// Validate and store the incoming secret name
	secretName := strings.TrimSpace(secretNameIn)
	if secretName == "" {
		return nil, fmt.Errorf("vault secret name provided to NewSealHandler cannot be empty")
	}
	componentLogger.Debug("Using Vault secret name", zap.String("secretName", secretName))

	hashiCorpVaultSealHandler, err := vault.NewHashiCorpVaultSealHandler(ctx, componentLogger, cfg.Vault.HashiCorp)
	if err != nil {
		// Log error and return it
		componentLogger.Error("Failed to create underlying HashiCorp Vault seal handler", zap.Error(err))
		return nil, fmt.Errorf("new hashicorp vault seal handler: %w", err)
	}
	componentLogger.Debug("Underlying HashiCorp Vault seal handler created")

	kuberConfig, err := rest.InClusterConfig()
	if err != nil {
		componentLogger.Error("Failed to get in-cluster Kubernetes config", zap.Error(err))
		return nil, fmt.Errorf("failed to get kubernetes config: %w", err)
	}
	clientset, err := kubernetes.NewForConfig(kuberConfig)
	if err != nil {
		componentLogger.Error("Failed to create Kubernetes clientset", zap.Error(err))
		return nil, fmt.Errorf("failed to create kubernetes clientset: %w", err)
	}
	componentLogger.Debug("Kubernetes clientset created")

	handler := &SealHandler{
		logger:           componentLogger,
		cfg:              cfg,
		vaultSealHandler: hashiCorpVaultSealHandler,
		kubeClientset:    clientset,
		secretName:       secretName, // <<< ASSIGN parameter to the struct field
	}

	componentLogger.Info("VaultSealHandler initialized successfully", zap.String("secretName", handler.secretName))
	return handler, nil
}

// initVault checks Vault initialization status, initializes if necessary,
// and ensures the keys secret state is consistent. Returns (bool, error).
// Uses s.secretName now.
func (s *SealHandler) initVault(ctx context.Context) (bool, error) {
	s.logger.Debug("Entering initVault")
	defer s.logger.Debug("Exiting initVault")
	s.logger.Debug("Checking Vault initialization status...")
	initRes, err := s.vaultSealHandler.TryInit(ctx)
	if err != nil {
		s.logger.Error("Failed checking vault init status", zap.Error(err))
		return false, fmt.Errorf("%w: %w", ErrVaultInitCheckFailed, err)
	}

	if initRes != nil {
		s.logger.Info("Vault was not initialized. Initialization performed, storing keys.")
		keysSecret := corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: s.secretName /* Use field */, Namespace: s.cfg.OpengovernanceNamespace}, Type: corev1.SecretTypeOpaque, StringData: make(map[string]string)}
		for i, key := range initRes.Keys {
			keysSecret.StringData[fmt.Sprintf("key-%d", i)] = key
		}
		keysSecret.StringData["root-token"] = initRes.RootToken

		secretsClient := s.kubeClientset.CoreV1().Secrets(s.cfg.OpengovernanceNamespace)
		s.logger.Info("Attempting to create Kubernetes secret for unseal keys", zap.String("secretName", s.secretName /* Use field */))
		_, err = secretsClient.Create(ctx, &keysSecret, metav1.CreateOptions{})
		if err != nil {
			if k8serrors.IsAlreadyExists(err) {
				s.logger.Warn("Vault keys secret already existed unexpectedly after init, attempting update.", zap.String("secretName", s.secretName /* Use field */))
				_, updateErr := secretsClient.Update(ctx, &keysSecret, metav1.UpdateOptions{})
				if updateErr != nil {
					s.logger.Error("Failed to update existing vault keys secret after initialization", zap.String("secretName", s.secretName), zap.Error(updateErr))
					return false, fmt.Errorf("%w: %w", ErrVaultSecretUpdateFailed, updateErr)
				}
				s.logger.Info("Successfully updated existing vault keys secret.", zap.String("secretName", s.secretName /* Use field */))
			} else {
				s.logger.Error("Failed to create vault keys secret after initialization", zap.String("secretName", s.secretName), zap.Error(err))
				return false, fmt.Errorf("%w: %w", ErrVaultSecretCreateFailed, err)
			}
		} else {
			s.logger.Info("Successfully created vault keys secret.", zap.String("secretName", s.secretName /* Use field */))
		}
		s.logger.Info("Vault initialization complete and keys stored.")
		return true, nil // New init, success
	} else {
		s.logger.Info("Vault already initialized. Verifying keys secret exists.")
		secretsClient := s.kubeClientset.CoreV1().Secrets(s.cfg.OpengovernanceNamespace)
		_, err := secretsClient.Get(ctx, s.secretName /* Use field */, metav1.GetOptions{})
		if err != nil {
			errMsg := ""
			var specificErr error
			if k8serrors.IsNotFound(err) {
				errMsg = fmt.Sprintf("inconsistent State: Vault initialized, but keys secret '%s' is missing in namespace '%s'", s.secretName /* Use field */, s.cfg.OpengovernanceNamespace)
				specificErr = ErrVaultKeysMissing
			} else {
				errMsg = fmt.Sprintf("failed to get vault keys secret '%s' for verification", s.secretName /* Use field */)
				specificErr = ErrVaultSecretGetFailed
			}
			s.logger.Error(errMsg, zap.Error(err))
			return false, fmt.Errorf("%w: %w", specificErr, err) // Wrap original k8s error
		}
		s.logger.Info("Vault initialized state and keys secret existence verified.", zap.String("secretName", s.secretName /* Use field */))
		return false, nil // No new init, success
	}
}

// vault/seal.go

// unsealChecker attempts initial unseal & setup, then runs periodic checks.
// Uses s.secretName now.
func (s *SealHandler) unsealChecker(ctx context.Context, isNewInit bool, unsealed chan<- struct{}) {
	s.logger.Debug("Starting unsealChecker goroutine", zap.Bool("isNewInit", isNewInit))
	signaledSuccessfully := false

	defer func() { // Defer function remains the same
		if r := recover(); r != nil {
			s.logger.Error("Panic recovered in unsealChecker periodic loop.", zap.Any("panicValue", r), zap.String("stacktrace", string(debug.Stack())))
		}
		if !signaledSuccessfully && unsealed != nil {
			s.logger.Warn("Closing unseal channel due to unsealChecker exiting early.")
			close(unsealed)
		}
		s.logger.Debug("unsealChecker goroutine finished.")
	}()

	// 1. Get Keys
	s.logger.Debug("Attempting to get unseal keys from Kubernetes secret", zap.String("secretName", s.secretName /* Use field */))
	secretsClient := s.kubeClientset.CoreV1().Secrets(s.cfg.OpengovernanceNamespace)
	keysSecret, err := secretsClient.Get(ctx, s.secretName /* Use field */, metav1.GetOptions{})
	if err != nil {
		s.logger.Error("CRITICAL: Failed to get vault unseal keys secret at start. Cannot proceed.", zap.String("secretName", s.secretName), zap.Error(err))
		return // Exit goroutine, defer closes channel
	}
	s.logger.Debug("Successfully retrieved unseal keys secret", zap.String("secretName", s.secretName /* Use field */))

	keys := make([]string, 0, len(keysSecret.Data))
	var rootTokenBytes []byte
	var rootTokenFound bool
	for k, v := range keysSecret.Data {
		if k == "root-token" {
			rootTokenBytes = v
			rootTokenFound = true
			continue
		}
		keys = append(keys, string(v))
	}
	s.logger.Debug("Extracted unseal keys from secret", zap.Int("keyCount", len(keys)), zap.Bool("rootTokenFound", rootTokenFound))
	if len(keys) == 0 {
		s.logger.Error("CRITICAL: No unseal keys found in secret data. Cannot proceed.", zap.String("secretName", s.secretName /* Use field */))
		return // Exit goroutine, defer closes channel
	}

	// 2. Initial Unseal & Auth Setup
	initialAttemptDone := false
	if unsealed != nil {
		s.logger.Debug("Performing initial unseal attempt...")
		err = s.vaultSealHandler.TryUnseal(ctx, keys)
		if err != nil {
			s.logger.Error("Initial unseal attempt failed, will retry periodically.", zap.Error(err))
		} else {
			s.logger.Info("Initial unseal successful")
			initialAttemptDone = true
			if rootTokenFound && len(rootTokenBytes) > 0 {
				s.logger.Info("Attempting to setup Kubernetes auth in Vault...")
				kubeAuthErr := s.vaultSealHandler.SetupKuberAuth(ctx, string(rootTokenBytes))
				if kubeAuthErr != nil {
					s.logger.Error("Failed to setup Kubernetes auth", zap.Error(kubeAuthErr))
				} else {
					s.logger.Info("Kubernetes auth setup successful or already configured.")
				}
			} else {
				s.logger.Warn("Root token not found/empty in secret, skipping Kubernetes auth setup.", zap.String("secretName", s.secretName /* Use field */))
			}
			s.logger.Debug("Signaling initial unseal complete.")
			select {
			case unsealed <- struct{}{}:
				signaledSuccessfully = true
				close(unsealed)
			case <-ctx.Done():
				s.logger.Warn("Context cancelled before initial unseal signal could be sent.")
				close(unsealed)
			}
			unsealed = nil // Prevent double close in defer
		}
	} else {
		s.logger.Warn("unsealChecker called with nil channel.")
	}

	// 3. Periodic Loop
	if unsealed == nil && !initialAttemptDone {
		return
	} // Exit if channel was nil initially
	ticker := time.NewTicker(unsealCheckInterval)
	defer ticker.Stop()
	if initialAttemptDone {
		s.logger.Info("Starting periodic unseal check loop", zap.Duration("interval", unsealCheckInterval))
	} else {
		s.logger.Warn("Initial unseal failed, starting periodic retry loop", zap.Duration("interval", unsealCheckInterval))
	}
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Context cancelled, stopping unseal checker.", zap.Error(ctx.Err()))
			return
		case tickTime := <-ticker.C:
			s.logger.Debug("Periodic unseal check", zap.Time("tickTime", tickTime))
			if err = s.vaultSealHandler.TryUnseal(ctx, keys); err != nil {
				s.logger.Warn("Periodic unseal attempt failed", zap.Error(err))
			} else {
				s.logger.Debug("Periodic unseal check successful/already unsealed")
			}
		}
	}

} // --- END unsealChecker ---
// Start initiates the Vault initialization and unsealing process.
// It blocks until the initial unseal attempt completes successfully or a timeout occurs.
// Returns an error if init fails or if the initial unseal fails/times out.
func (s *SealHandler) Start(ctx context.Context) error { // <<< Changed signature
	s.logger.Debug("Entering VaultSealHandler Start")

	s.logger.Info("Performing Vault initialization check...")
	isNewInit, initErr := s.initVault(ctx)
	if initErr != nil {
		// *** CHANGED: Return error instead of panic ***
		return fmt.Errorf("vault initialization check failed: %w", initErr)
	}
	s.logger.Info("Vault initialization check completed.", zap.Bool("isNewInit", isNewInit))

	unsealChan := make(chan struct{})
	s.logger.Info("Starting Vault unseal checker background goroutine...")
	go s.unsealChecker(ctx, isNewInit, unsealChan)

	s.logger.Info("Waiting for initial Vault unseal signal...", zap.Duration("timeout", initialUnsealTimeout))
	select {
	case _, ok := <-unsealChan:
		if ok {
			// Received signal (struct{}{}), channel still open or just closed by sender.
			s.logger.Info("Initial unseal signal received. Vault should be ready.")
			// Success
		} else {
			// Channel was closed without sending a value (checker exited early).
			// *** CHANGED: Return specific error ***
			s.logger.Error("Vault unseal checker exited prematurely before signaling success.")
			return ErrUnsealCheckerExited // Return specific error
		}
	case <-time.After(initialUnsealTimeout):
		// *** CHANGED: Return specific error ***
		s.logger.Error("Timeout waiting for initial Vault unseal signal.", zap.Duration("timeout", initialUnsealTimeout))
		return fmt.Errorf("%w (%v)", ErrUnsealTimeout, initialUnsealTimeout)
	case <-ctx.Done():
		// *** CHANGED: Return context error ***
		s.logger.Error("Context cancelled while waiting for initial Vault unseal signal.", zap.Error(ctx.Err()))
		return fmt.Errorf("context cancelled during vault unseal wait: %w", ctx.Err())
	}

	s.logger.Info("VaultSealHandler Start completed successfully.")
	return nil // <<< Return nil on success
}
