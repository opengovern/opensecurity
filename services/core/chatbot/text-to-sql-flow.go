package chatbot

import (
	"bytes"
	"context"
	"fmt"
	"github.com/opengovern/opensecurity/services/core/db/models"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
	"time"

	"gopkg.in/yaml.v3"
)

// QueryAttempt represents an attempt to generate and validate an SQL query.
type QueryAttempt struct {
	Query string `json:"query"`
	Error string `json:"error"`
}

type RequestData struct {
	Question         string         `json:"question"`
	PreviousAttempts []QueryAttempt `json:"previous_attempts"`
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
func NewTextToSQLFlow() (*TextToSQLFlow, error) {
	appConfig, err := NewAppConfig()
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

func extractSQLFromResponse(responseText string) string {
	// Example 1: Look for ```sql ... ``` block
	sqlBlockStart := "```sql"
	sqlBlockEnd := "```"
	startIdx := strings.Index(responseText, sqlBlockStart)
	if startIdx != -1 {
		endIdx := strings.Index(responseText[startIdx+len(sqlBlockStart):], sqlBlockEnd)
		if endIdx != -1 {
			return strings.TrimSpace(responseText[startIdx+len(sqlBlockStart) : startIdx+len(sqlBlockStart)+endIdx])
		}
	}

	// Example 2: Look for first SELECT statement (very basic)
	upperResponse := strings.ToUpper(responseText)
	selectIdx := strings.Index(upperResponse, "SELECT")
	if selectIdx != -1 {
		// Find the end (e.g., semicolon or end of string) - this is naive
		endIdx := strings.Index(responseText[selectIdx:], ";")
		if endIdx != -1 {
			return strings.TrimSpace(responseText[selectIdx : selectIdx+endIdx+1])
		}
		return strings.TrimSpace(responseText[selectIdx:]) // Return rest of string if no semicolon
	}

	return strings.TrimSpace(responseText)
}

func (f *TextToSQLFlow) RunInference(ctx context.Context, chat *models.Chat, data RequestData, agentInput *string) (agent string, finalQuery string, err error) {
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
			return "", "", fmt.Errorf("failed during IAM verification: %w", err)
		}

		// Normalize: strip whitespace, remove asterisks, trailing colons, and convert to lowercase.
		verifierResultClean := strings.ToLower(strings.TrimSpace(strings.ReplaceAll(strings.TrimSuffix(verifierResult, ":"), "*", "")))

		if _, exists := f.mappingData[verifierResultClean]; !exists {
			// Explicitly check if the detected key exists in the map
			return "", "", fmt.Errorf("request is not supported. Detected agent '%s' (from '%s') is not configured in mapping", verifierResultClean, verifierResult)
		}
		agent = verifierResultClean
		log.Printf("Debug: Determined agent from classification: %s", agent)
	}

	// --- 2. Retrieve Agent Configuration ---
	agentConfigData, ok := f.mappingData[agent]
	if !ok {
		return "", "", fmt.Errorf("configuration for agent '%s' not found or not a map in mapping data", agent)
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
			return "", "", fmt.Errorf("schema file '%s' not found: %w", schemaFilePath, err)
		}
		schemaBytes, err := os.ReadFile(schemaFilePath)
		if err != nil {
			return "", "", fmt.Errorf("failed to read schema file '%s': %w", schemaFilePath, err)
		}
		schemaBuilder.Write(schemaBytes)
		schemaBuilder.WriteString("\n") // Add newline between files
	}
	schemaText := schemaBuilder.String()

	// --- 4. Load and Prepare Prompt Template ---
	if _, err := os.Stat(promptFile); os.IsNotExist(err) {
		return "", "", fmt.Errorf("prompt file '%s' not found: %w", promptFile, err)
	}
	promptYamlBytes, err := os.ReadFile(promptFile)
	if err != nil {
		return "", "", fmt.Errorf("failed to read prompt file '%s': %w", promptFile, err)
	}

	var promptData PromptData
	err = yaml.Unmarshal(promptYamlBytes, &promptData)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse prompt YAML '%s': %w", promptFile, err)
	}

	// --- 5. Build Template Data ---
	templateRenderData := make(map[string]any)

	templateRenderData["question"] = data.Question

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

		tmpl, err := template.New("prompt").Parse(templateContent)
		if err != nil {
			return "", "", fmt.Errorf("failed to parse prompt template for role %s: %w", role, err)
		}

		var renderedContent bytes.Buffer
		err = tmpl.Execute(&renderedContent, templateRenderData)
		if err != nil {
			return "", "", fmt.Errorf("failed to render prompt template for role %s: %w", role, err)
		}

		messages = append(messages, ChatMessage{Role: role, Content: renderedContent.String()})
	}

	responseText, err := f.llmClient.ChatCompletion(ctx, modelName, messages, 800, 0.0) // Use defaults or get from config
	if err != nil {
		return agent, "", fmt.Errorf("text-to-SQL LLM call failed: %w", err)
	}
	log.Printf("Debug: LLM raw response:\n%s", responseText)

	if responseText == "" {
		return agent, "", fmt.Errorf("no response from text-to-SQL model")
	}
	if strings.Contains(strings.ToUpper(responseText), "ERROR") {
		return agent, "", fmt.Errorf("model returned 'ERROR': %s", responseText)
	}

	finalQuery = extractSQLFromResponse(responseText)
	if finalQuery == "" {
		return agent, "", fmt.Errorf("no valid SQL found in model response")
	}

	log.Printf("Info: Final SQL:\n%s", finalQuery)
	return agent, finalQuery, nil
}
