package chatbot

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/flosch/pongo2/v6"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

// QueryAttempt represents an attempt to generate and validate an SQL query.
type QueryAttempt struct {
	Query string `json:"query"`
	Error string `json:"error"`
}

type RequestData struct {
	Question                  string         `json:"question"`
	PreviousAttempts          []QueryAttempt `json:"previous_attempts"`
	InClarificationState      bool           `json:"in_clarification_state"`
	ClarificationQuestions    []string       `json:"clarification_questions"`
	UserClarificationResponse string         `json:"user_clarification_response"`
}

// TextToSQLFlow converts natural language questions to SQL queries.
type TextToSQLFlow struct {
	llmClient     *LLMClient
	mappingData   MappingData
	queryAttempts []*QueryAttempt
	mu            sync.RWMutex
}

type PromptData struct {
	Prompts []struct {
		Role    string `yaml:"role"`
		Content string `yaml:"content"`
	} `yaml:"prompts"`
}

// NewTextToSQLFlow creates a new TextToSQLFlow instance.
// baseDir is the root directory for resolving relative paths in mapping_data.
func NewTextToSQLFlow(hfToken string) (*TextToSQLFlow, error) {
	appConfig, err := NewAppConfig(hfToken)
	if err != nil {
		return nil, err
	}

	llmClient, err := NewLLMClient(appConfig.HfToken, appConfig.GetProvider())

	return &TextToSQLFlow{
		llmClient:     llmClient,
		mappingData:   appConfig.MappingData,
		queryAttempts: make([]*QueryAttempt, 0),
		mu:            sync.RWMutex{},
	}, nil
}

// AddQueryAttempt adds a new query attempt to the list (thread-safe).
func (f *TextToSQLFlow) AddQueryAttempt(query string, errorMsg string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	attempt := &QueryAttempt{Query: query, Error: errorMsg}
	f.queryAttempts = append(f.queryAttempts, attempt)
	log.Printf("Debug: Added query attempt: %+v", attempt)
}

// ClearQueryAttempts clears all stored query attempts (thread-safe).
func (f *TextToSQLFlow) ClearQueryAttempts() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.queryAttempts = make([]*QueryAttempt, 0)
	log.Println("Debug: Cleared all query attempts.")
}

// parseLLMResponse attempts to parse the LLM's text response into a structured InferenceResult.
// It handles different JSON formats (SUCCESS, CLARIFICATION_NEEDED, ERROR)
// and falls back to checking for raw SQL if JSON parsing fails initially.
func parseLLMResponse(responseText string) *InferenceResult {
	// Clean up potential markdown fences and whitespace
	trimmedResponse := strings.TrimSpace(responseText)
	trimmedResponse = strings.TrimPrefix(trimmedResponse, "```json")
	trimmedResponse = strings.TrimPrefix(trimmedResponse, "```")
	trimmedResponse = strings.TrimSuffix(trimmedResponse, "```")
	trimmedResponse = strings.TrimSpace(trimmedResponse) // Trim again after removing fences

	if trimmedResponse == "" {
		return &InferenceResult{
			Type:        ResultTypeMalformed,
			Reason:      "LLM returned an empty response.",
			RawResponse: responseText,
		}
	}

	// Attempt to parse as JSON to determine the result type
	var base BaseResponse
	err := json.Unmarshal([]byte(trimmedResponse), &base)

	if err != nil {
		// JSON parsing failed initially. Check if it's just raw SQL.
		log.Printf("Debug: Initial JSON parse failed: %v. Checking for raw SQL fallback.", err)
		potentialSQL := extractSQLFallback(trimmedResponse) // Use a helper for old logic
		if potentialSQL != "" {
			log.Printf("Debug: Found raw SQL using fallback logic.")
			return &InferenceResult{
				Type:                  ResultTypeRawSQL, // Indicate it was found via fallback
				Query:                 potentialSQL,
				RawResponse:           responseText,                       // Store original response
				PrimaryInterpretation: "Raw SQL extracted from response.", // Add default interpretation
			}
		}
		// If not valid JSON and not raw SQL, consider it malformed.
		log.Printf("ERROR: Response is not valid JSON and does not appear to contain raw SQL.")
		return &InferenceResult{
			Type:        ResultTypeMalformed,
			Reason:      fmt.Sprintf("Response is not valid JSON and does not contain detectable SQL: %v", err),
			RawResponse: responseText,
		}
	}

	// JSON parsing succeeded, now parse into the specific structure based on 'result' field
	switch base.Result {
	case ResultTypeSuccess:
		var successResp SuccessResponse
		if err := json.Unmarshal([]byte(trimmedResponse), &successResp); err != nil {
			return &InferenceResult{Type: ResultTypeMalformed, Reason: fmt.Sprintf("failed to parse SUCCESS response: %v", err), RawResponse: responseText}
		}
		// Validate essential fields
		if successResp.Query == "" {
			return &InferenceResult{Type: ResultTypeMalformed, Reason: "SUCCESS response missing 'query' field", RawResponse: responseText}
		}
		return &InferenceResult{
			Type:                      ResultTypeSuccess,
			Query:                     successResp.Query,
			PrimaryInterpretation:     successResp.PrimaryInterpretation,
			AdditionalInterpretations: successResp.AdditionalInterpretations,
		}

	case ResultTypeClarificationNeeded:
		var clarResp ClarificationResponse
		if err := json.Unmarshal([]byte(trimmedResponse), &clarResp); err != nil {
			return &InferenceResult{Type: ResultTypeMalformed, Reason: fmt.Sprintf("failed to parse CLARIFICATION_NEEDED response: %v", err), RawResponse: responseText}
		}
		// Validate essential fields
		if len(clarResp.ClarifyingQuestions) == 0 {
			return &InferenceResult{Type: ResultTypeMalformed, Reason: "CLARIFICATION_NEEDED response missing 'clarifying_questions'", RawResponse: responseText}
		}
		return &InferenceResult{
			Type:                ResultTypeClarificationNeeded,
			ClarifyingQuestions: clarResp.ClarifyingQuestions,
		}

	case ResultTypeError:
		var errResp ErrorResponse
		if err := json.Unmarshal([]byte(trimmedResponse), &errResp); err != nil {
			return &InferenceResult{Type: ResultTypeMalformed, Reason: fmt.Sprintf("failed to parse ERROR response: %v", err), RawResponse: responseText}
		}
		// Validate essential fields
		if errResp.Reason == "" {
			return &InferenceResult{Type: ResultTypeMalformed, Reason: "ERROR response missing 'reason' field", RawResponse: responseText}
		}
		return &InferenceResult{
			Type:   ResultTypeError,
			Reason: errResp.Reason,
		}

	default:
		// Unknown 'result' value
		return &InferenceResult{
			Type:        ResultTypeMalformed,
			Reason:      fmt.Sprintf("LLM response contained unknown result type: '%s'", base.Result),
			RawResponse: responseText,
		}
	}
}

// extractSQLFallback contains the previous logic for finding SQL in non-JSON responses.
func extractSQLFallback(responseText string) string {
	// Example 1: Look for ```sql ... ``` block
	sqlBlockStart := "```sql"
	sqlBlockEnd := "```"
	startIdx := strings.Index(responseText, sqlBlockStart)
	if startIdx != -1 {
		endIdx := strings.Index(responseText[startIdx+len(sqlBlockStart):], sqlBlockEnd)
		if endIdx != -1 {
			return strings.TrimSpace(responseText[startIdx+len(sqlBlockStart) : startIdx+len(sqlBlockStart)+endIdx])
		}
		// Handle case where closing fence is missing but opening exists
		return strings.TrimSpace(responseText[startIdx+len(sqlBlockStart):])
	}

	// Example 2: Look for first SELECT statement (very basic)
	// Only consider SELECT if it appears early in the string, otherwise might be part of explanation
	upperResponse := strings.ToUpper(strings.TrimSpace(responseText))
	if strings.HasPrefix(upperResponse, "SELECT") {
		// Heuristic: If it starts with SELECT, assume it's SQL. This is fragile.
		// Find the end (e.g., semicolon or end of string) - this is naive
		endIdx := strings.Index(responseText, ";")
		if endIdx != -1 {
			return strings.TrimSpace(responseText[:endIdx+1])
		}
		return strings.TrimSpace(responseText) // Return whole string if no semicolon
	}

	return "" // Return empty if no SQL found via fallback
}

func (f *TextToSQLFlow) RunInference(ctx context.Context, data RequestData, agentInput *string) (agent string, finalResponse *InferenceResult, err error) {
	question := strings.TrimSpace(data.Question)

	log.Println("Debug: === run_inference called ===")
	log.Printf("Debug: Question: %s", question)

	// --- 1. Determine Agent (Domain) ---
	if agentInput != nil && *agentInput != "" {
		agent = strings.ToLower(strings.TrimSpace(*agentInput))
		log.Printf("Debug: Using provided agent: %s", agent)
	} else {
		verifierPromptFile := ClassifyPromptPath
		defaultModelName := "Qwen/Qwen2.5-72B-Instruct"

		var verifierResult string
		verifierResult, err = f.llmClient.Verify(ctx, question, verifierPromptFile, defaultModelName)
		if err != nil {
			// Decide if this error should halt execution or just default agent
			return "", nil, fmt.Errorf("failed during IAM verification: %w", err)
		}

		// Normalize: strip whitespace, remove asterisks, trailing colons, and convert to lowercase.
		verifierResultClean := strings.ToLower(strings.TrimSpace(strings.ReplaceAll(strings.TrimSuffix(verifierResult, ":"), "*", "")))

		if _, exists := f.mappingData[verifierResultClean]; !exists {
			// Explicitly check if the detected key exists in the map
			return "", nil, fmt.Errorf("request is not supported. Detected agent '%s' (from '%s') is not configured in mapping", verifierResultClean, verifierResult)
		}
		agent = verifierResultClean
		log.Printf("Debug: Determined agent from classification: %s", agent)
	}

	// --- 2. Retrieve Agent Configuration ---
	agentConfigData, ok := f.mappingData[agent]
	if !ok {
		return "", nil, fmt.Errorf("configuration for agent '%s' not found or not a map in mapping data", agent)
	}

	// Safely get config values with type assertions and defaults
	agentSpecificConfig := agentConfigData.AgentSpecificConfig
	modelName := agentSpecificConfig.PrimaryModel
	if modelName == "" {
		modelName = "Qwen/Qwen2.5-72B-Instruct" // Default model
	}

	promptFile := filepath.Join(PromptsPath, agentConfigData.PromptTemplateFile)

	log.Printf("Debug: Using agent=%s, model_name=%s, schema_files=%v, prompt_file=%s", agent, modelName, agentConfigData.SQLSchemaFiles, promptFile)

	// --- 3. Load Schema ---
	var schemaBuilder strings.Builder
	for _, schemaFileRel := range agentConfigData.SQLSchemaFiles {
		schemaFilePath := filepath.Join(SchemasPath, schemaFileRel)
		if _, err := os.Stat(schemaFilePath); os.IsNotExist(err) {
			return "", nil, fmt.Errorf("schema file '%s' not found: %w", schemaFilePath, err)
		}
		schemaBytes, err := os.ReadFile(schemaFilePath)
		if err != nil {
			return "", nil, fmt.Errorf("failed to read schema file '%s': %w", schemaFilePath, err)
		}
		schemaBuilder.Write(schemaBytes)
		schemaBuilder.WriteString("\n") // Add newline between files
	}
	schemaText := schemaBuilder.String()

	// --- 4. Load and Prepare Prompt Template ---
	if _, err := os.Stat(promptFile); os.IsNotExist(err) {
		return "", nil, fmt.Errorf("prompt file '%s' not found: %w", promptFile, err)
	}
	promptYamlBytes, err := os.ReadFile(promptFile)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read prompt file '%s': %w", promptFile, err)
	}

	var promptData PromptData
	err = yaml.Unmarshal(promptYamlBytes, &promptData)
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse prompt YAML '%s': %w", promptFile, err)
	}

	// --- 5. Build Template Data ---
	templateRenderData := pongo2.Context{}

	if data.PreviousAttempts == nil || len(data.PreviousAttempts) == 0 {
		f.mu.RLock()
		for _, qa := range f.queryAttempts {
			if qa != nil {
				templateRenderData["previous_attempts"] = *qa
			}
		}
		f.mu.RUnlock()
	} else {
		templateRenderData["previous_attempts"] = data.PreviousAttempts
	}

	templateRenderData["user_clarification_response"] = data.UserClarificationResponse
	templateRenderData["clarifying_questions"] = data.ClarificationQuestions
	templateRenderData["in_clarification_state"] = data.InClarificationState
	templateRenderData["domain_topic"] = agent
	templateRenderData["original_question"] = question
	templateRenderData["schema_text"] = schemaText
	templateRenderData["sql_engine"] = "PostgreSQL"
	templateRenderData["today_time"] = time.Now().UTC().Format("2006-01-02 15:04:05")

	messages := make([]ChatMessage, 0, len(promptData.Prompts))
	for _, prompt := range promptData.Prompts {
		role := prompt.Role
		if role == "" {
			role = "user"
		}

		templateContent := prompt.Content
		if role == "user" {
			if override := os.Getenv("USER_TEXT2SQL_PROMPT_TEMPLATE"); override != "" {
				templateContent = override
			}
		} else if role == "system" {
			if override := os.Getenv("SYS_TEXT2SQL_PROMPT_TEMPLATE"); override != "" {
				templateContent = override
			}
		}

		// Compile the template string using pongo2
		tmpl, err := pongo2.FromString(templateContent)
		if err != nil {
			return "", nil, fmt.Errorf("failed to compile pongo2 template for role %s: %w", role, err)
		}

		// Execute the template with the context data
		renderedContent, err := tmpl.Execute(templateRenderData)
		if err != nil {
			return "", nil, fmt.Errorf("failed to render pongo2 template for role %s: %w", role, err)
		}

		messages = append(messages, ChatMessage{Role: role, Content: renderedContent})
	}

	responseText, err := f.llmClient.ChatCompletion(ctx, modelName, messages, 800, 0.0) // Use defaults or get from config
	if err != nil {
		return agent, nil, fmt.Errorf("text-to-SQL LLM call failed: %w", err)
	}
	log.Printf("Debug: LLM raw response:\n%s", responseText)

	if responseText == "" {
		return agent, nil, fmt.Errorf("no response from text-to-SQL model")
	}
	if strings.Contains(strings.ToUpper(responseText), "ERROR") {
		return agent, nil, fmt.Errorf("model returned 'ERROR': %s", responseText)
	}

	finalResponse = parseLLMResponse(responseText)
	if finalResponse == nil {
		return agent, nil, fmt.Errorf("no valid response found")
	}

	log.Printf("Info: Final SQL:\n%s", *finalResponse)
	return agent, finalResponse, nil
}
