package chatbot

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

// AppConfig holds application configuration loaded from various sources.
type AppConfig struct {
	MappingData       MappingData
	HfToken           string
	Auth0Domain       string
	Auth0ClientID     string
	Auth0ClientSecret string
}

type AgentSubConfig struct {
	PrimaryModel    string `yaml:"primary_model"`
	PrimaryProvider string `yaml:"primary_provider"`
}

type AgentConfig struct {
	Name                     string         `yaml:"name"`
	Description              string         `yaml:"description"`
	WelcomeMessage           string         `yaml:"welcome_message"`
	SampleQuestions          []string       `yaml:"sample_questions"`
	Availability             string         `yaml:"availability"`
	PromptTemplateFile       string         `yaml:"prompt_template_file"`
	QueryVerificationRetries *int           `yaml:"query_verification_retries,omitempty"`
	SeekClarification        *bool          `yaml:"seek_clarification,omitempty"`
	Domains                  []string       `yaml:"domains"`
	SQLSchemaFiles           []string       `yaml:"sql_schema_files"`
	AgentSpecificConfig      AgentSubConfig `yaml:"agent_config"`
}

type MappingData map[string]AgentConfig

// NewAppConfig parses args, loads env vars and mapping file.
// mappingFilePath is the path to mapping.yaml.
func NewAppConfig() (*AppConfig, error) {
	cfg := &AppConfig{}

	mappingFilePath := MappingPath

	if _, err := os.Stat(mappingFilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("mapping file '%s' not found: %w", mappingFilePath, err)
	} else if err != nil {
		return nil, fmt.Errorf("error checking mapping file '%s': %w", mappingFilePath, err)
	}

	yamlFile, err := os.ReadFile(mappingFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read mapping file '%s': %w", mappingFilePath, err)
	}

	err = yaml.Unmarshal(yamlFile, &cfg.MappingData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML mapping file '%s': %w", mappingFilePath, err)
	}

	// --- Load Environment Variables ---
	cfg.HfToken = os.Getenv("HF_API_TOKEN")
	if cfg.HfToken == "" {
		// Consider returning error instead of panic/fatal for library code
		return nil, fmt.Errorf("please set HF_API_TOKEN environment variable")
	}

	cfg.Auth0Domain = os.Getenv("AUTH0_DOMAIN")
	cfg.Auth0ClientID = os.Getenv("AUTH0_CLIENT_ID")
	cfg.Auth0ClientSecret = os.Getenv("AUTH0_CLIENT_SECRET")

	// Optionally set BaseDir based on mappingFilePath or another env var
	// cfg.BaseDir = filepath.Dir(mappingFilePath) // Example

	log.Println("AppConfig initialized successfully.")
	return cfg, nil
}

// GetProvider retrieves the provider from the mapping data.
func (c *AppConfig) GetProvider() string {
	return "together"
}
