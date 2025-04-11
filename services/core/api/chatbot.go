package api

import "github.com/opengovern/opensecurity/services/core/chatbot"

// QueryAttempt represents an attempt to generate and validate an SQL query.
type QueryAttempt struct {
	Query string `json:"query"`
	Error string `json:"error"`
}
type GenerateQueryRequest struct {
	Question         string         `json:"question"`
	PreviousAttempts []QueryAttempt `json:"previous_attempts"`
	Agent            *string        `json:"agent,omitempty"`
	RetryCount       *int           `json:"retry_count,omitempty"`
}

type GenerateQueryResponse struct {
	Result chatbot.InferenceResult `json:"result"`
	Agent  string                  `json:"agent"`
}

type AttemptResult struct {
	Result   chatbot.InferenceResult `json:"result"`
	Agent    string                  `json:"agent"`
	RunError *string                 `json:"run_error,omitempty"`
}

type GenerateQueryAndRunResponse struct {
	RunResult       RunQueryResponse `json:"result"`
	AttemptsResults []AttemptResult  `json:"attempts_results"`
}

type ConfigureChatbotSecretRequest struct {
	Key    string `json:"key"`
	Secret string `json:"secret"`
}
