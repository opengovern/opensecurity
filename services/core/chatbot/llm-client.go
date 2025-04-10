package chatbot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
)

const (
	// DefaultHuggingFaceAPIEndpoint is assumed based on common patterns.
	// Verify if this is correct for your specific provider/model.
	DefaultHuggingFaceAPIEndpoint = "https://api-inference.huggingface.co/v1/chat/completions"
	defaultMaxTokens              = 800
	defaultTemperature            = 0.0
)

// PromptConfig defines the top-level structure of the prompt YAML file.
type PromptConfig struct {
	Prompts []ChatMessage `yaml:"prompts"`
}

// ChatMessage corresponds to one message in the chat history.
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionRequest is the structure for the API request body.
type ChatCompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float64       `json:"temperature"`
	Stream      bool          `json:"stream"`
}

// ChatCompletionResponseChoice is a single choice in the API response.
type ChatCompletionResponseChoice struct {
	Message ChatMessage `json:"message"`
	// Add other fields like index, finish_reason if needed based on actual API response
}

// ChatCompletionResponse is the structure for the API response body.
type ChatCompletionResponse struct {
	Choices []ChatCompletionResponseChoice `json:"choices"`
	// Add other fields like id, object, created, model, usage if needed
}

// LLMClient wraps the HTTP client and configuration for Hugging Face API calls.
type LLMClient struct {
	httpClient  *http.Client
	hfToken     string
	apiEndpoint string
}

// NewLLMClient creates and initializes a new LLMClient.
func NewLLMClient(hfToken string, provider string) (*LLMClient, error) {
	if hfToken == "" {
		return nil, fmt.Errorf("Hugging Face token (hf_token) cannot be empty")
	}

	debugModeStr := os.Getenv("DEBUG_MODE")
	debugMode := false
	if b, err := strconv.ParseBool(debugModeStr); err == nil {
		debugMode = b
	} else if strings.ToLower(debugModeStr) == "1" || strings.ToLower(debugModeStr) == "true" || strings.ToLower(debugModeStr) == "yes" {
		debugMode = true
	}

	skipPromptValue := os.Getenv("SKIP_PROMPT_PRECHECK")
	if skipPromptValue == "" {
		skipPromptValue = "not set"
	}

	log.Printf("Initializing LLMClient with provider='%s'.", provider)
	log.Printf("DEBUG_MODE=%t, SKIP_PROMPT_PRECHECK=%s", debugMode, skipPromptValue)
	log.Printf("All logs will be at INFO level (using standard logger) to show maximum detail, ignoring debug mode.")

	// You can customize the http.Client, e.g., set timeouts
	httpClient := &http.Client{
		Timeout: 60 * time.Second, // Example timeout
	}

	return &LLMClient{
		httpClient:  httpClient,
		hfToken:     hfToken,
		apiEndpoint: DefaultHuggingFaceAPIEndpoint, // Use default, can be made configurable
	}, nil
}

// ChatCompletion makes a remote chat completion call with retry logic.
func (c *LLMClient) ChatCompletion(ctx context.Context, modelName string, messages []ChatMessage, maxTokens int, temperature float64) (string, error) {

	// Use defaults if zero values are passed
	if maxTokens <= 0 {
		maxTokens = defaultMaxTokens
	}
	// Temperature 0.0 is a valid value, no default check needed unless explicitly desired

	// --- Logging ---
	log.Printf("Requesting chat completion: model='%s' (max_tokens=%d, temperature=%.2f)",
		modelName, maxTokens, temperature)
	// Log full prompt messages (consider marshaling for cleaner output if needed)
	log.Printf("Full prompt messages: %+v", messages) // %+v provides more detail for structs

	// --- Prepare Request ---
	requestPayload := ChatCompletionRequest{
		Model:       modelName,
		Messages:    messages,
		MaxTokens:   maxTokens,
		Temperature: temperature,
		Stream:      false, // Explicitly false as in the Python example
	}

	requestBody, err := json.Marshal(requestPayload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	// --- Retry Logic ---
	var responseText string
	operation := func() error {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiEndpoint, bytes.NewReader(requestBody))
		if err != nil {
			// Don't retry on request creation errors, they are likely programming errors
			return backoff.Permanent(fmt.Errorf("failed to create request: %w", err))
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+c.hfToken)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			// Retry on connection errors or temporary network issues
			log.Printf("WARN: Attempt failed, connection error: %v. Retrying...", err)
			return fmt.Errorf("http client error: %w", err) // Retryable error
		}
		defer resp.Body.Close()

		respBodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			// Harder to say if this is retryable, could be connection cut midway. Let's retry.
			log.Printf("WARN: Attempt failed, failed to read response body: %v. Retrying...", err)
			return fmt.Errorf("failed to read response body: %w", err)
		}

		// Check for non-successful HTTP status codes
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			errMessage := fmt.Sprintf("api request failed: status code %d, body: %s", resp.StatusCode, string(respBodyBytes))
			// Decide if the error is permanent (e.g., 4xx) or potentially retryable (e.g., 5xx)
			if resp.StatusCode >= 400 && resp.StatusCode < 500 {
				log.Printf("ERROR: Attempt failed with client error %d. Not retrying.", resp.StatusCode)
				return backoff.Permanent(fmt.Errorf(errMessage)) // Do not retry 4xx errors
			}
			log.Printf("WARN: Attempt failed with server error %d. Retrying...", resp.StatusCode)
			return fmt.Errorf(errMessage) // Retry 5xx errors or others
		}

		// --- Process successful response ---
		var completionResponse ChatCompletionResponse
		err = json.Unmarshal(respBodyBytes, &completionResponse)
		if err != nil {
			// If we got a 2xx response but can't parse it, it's unlikely retrying will help
			log.Printf("ERROR: Failed to unmarshal successful response body: %v. Body: %s", err, string(respBodyBytes))
			return backoff.Permanent(fmt.Errorf("failed to unmarshal response body: %w", err))
		}

		// Extract the response text
		if len(completionResponse.Choices) > 0 {
			// Assuming we always care about the first choice like in the Python code
			responseText = completionResponse.Choices[0].Message.Content
			log.Printf("Received full response text: %s", responseText)
			return nil // Success, stop retrying
		}

		log.Printf("No valid message content in completion response choices.")
		responseText = ""
		return nil // Success (API call worked), but no content found

	}

	// Exponential backoff configuration (similar to tenacity defaults)
	b := backoff.NewExponentialBackOff()
	b.InitialInterval = 1 * time.Second // Adjusted from Python's min=2 for slightly faster first retry
	b.MaxInterval = 10 * time.Second
	b.Multiplier = 1.5   // Exponential factor (tenacity default is 2, 1.5 is less aggressive)
	b.MaxElapsedTime = 0 // No total time limit, rely on attempts
	// Use backoff.WithMaxRetries for attempt limit
	retryWithAttempts := backoff.WithMaxRetries(b, 5) // 5 attempts total (1 initial + 4 retries)

	// Run the operation with retries
	err = backoff.Retry(operation, retryWithAttempts)
	if err != nil {
		log.Printf("ERROR: Chat completion failed after multiple retries: %v", err)
		return "", fmt.Errorf("chat completion failed after retries: %w", err)
	}

	return responseText, nil
}

func (v *LLMClient) Verify(ctx context.Context, question string, promptFile string, modelName string) (string, error) {
	// Check if prompt file exists
	if _, err := os.Stat(promptFile); os.IsNotExist(err) {
		return "", fmt.Errorf("verifier prompt file '%s' not found: %w", promptFile, err)
	} else if err != nil {
		return "", fmt.Errorf("error checking prompt file '%s': %w", promptFile, err)
	}

	// Read and parse YAML prompt file
	yamlFile, err := os.ReadFile(promptFile)
	if err != nil {
		return "", fmt.Errorf("failed to read prompt file '%s': %w", promptFile, err)
	}

	var promptConf PromptConfig
	err = yaml.Unmarshal(yamlFile, &promptConf)
	if err != nil {
		return "", fmt.Errorf("failed to parse YAML prompt file '%s': %w", promptFile, err)
	}

	// Build messages
	messages := make([]ChatMessage, 0, len(promptConf.Prompts))
	for _, prompt := range promptConf.Prompts {
		role := prompt.Role
		if role == "" {
			role = "system" // Default role
		}
		content := strings.ReplaceAll(prompt.Content, "{{ user_input }}", question)
		messages = append(messages, ChatMessage{Role: role, Content: content})
	}

	// Call LLM
	// Use low max_tokens and zero temperature for classification
	responseText, err := v.ChatCompletion(ctx, modelName, messages, 20, 0.0)
	if err != nil {
		// Log the underlying error but return a generic failure message or fallback
		log.Printf("WARN: IAMVerifier LLM call failed: %v. Falling back to OFF-TOPIC.", err)
		return "OFF-TOPIC", nil
	}

	responseText = strings.TrimSpace(responseText)

	if responseText == "" {
		log.Println("WARN: IAMVerifier received empty response from LLM. Falling back to OFF-TOPIC.")
		return "OFF-TOPIC", nil
	}

	return strings.ToUpper(responseText), nil
}
