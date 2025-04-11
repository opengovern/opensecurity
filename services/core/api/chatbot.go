package api

import (
	"github.com/opengovern/opensecurity/services/core/chatbot"
	"time"
)

// QueryAttempt represents an attempt to generate and validate an SQL query.
type QueryAttempt struct {
	Query string `json:"query"`
	Error string `json:"error"`
}
type GenerateQueryRequest struct {
	SessionId                 *string        `json:"session_id"`
	ChatId                    *string        `json:"chat_id"`
	Question                  string         `json:"question"`
	PreviousAttempts          []QueryAttempt `json:"previous_attempts"`
	Agent                     *string        `json:"agent,omitempty"`
	RetryCount                *int           `json:"retry_count,omitempty"`
	InClarificationState      bool           `json:"in_clarification_state"`
	ClarificationQuestions    []string       `json:"clarification_questions"`
	UserClarificationResponse string         `json:"user_clarification_response"`
}

type ClarificationQuestion struct {
	ClarificationId string `json:"clarification_id"`
	Question        string `json:"question"`
}

type Suggestion struct {
	SuggestionId string `json:"suggestion_id"`
	Suggestion   string `json:"suggestion"`
}

type InferenceResult struct {
	Type chatbot.ResultType `json:"type"`

	Query                     string       `json:"query,omitempty"`
	PrimaryInterpretation     Suggestion   `json:"primary_interpretation,omitempty"`
	AdditionalInterpretations []Suggestion `json:"additional_interpretations,omitempty"`

	ClarifyingQuestions []ClarificationQuestion `json:"clarifying_questions,omitempty"`

	Reason string `json:"reason,omitempty"`

	RawResponse string `json:"raw_response,omitempty"`
}

type GenerateQueryResponse struct {
	SessionId string          `json:"session_id"`
	ChatId    string          `json:"chat_id"`
	Result    InferenceResult `json:"result"`
	Agent     string          `json:"agent"`
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

type Session struct {
	ID      string `json:"id"`
	AgentId string `json:"agent_id"`
	Chats   []Chat `json:"chats"`
}

type Chat struct {
	ID                string         `json:"id"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	Question          string         `json:"question"`
	Query             *string        `json:"query,omitempty"`
	QueryError        *string        `json:"query_error,omitempty"`
	NeedClarification bool           `json:"need_clarification"`
	AssistantText     *string        `json:"assistant_text,omitempty"`
	TimeTaken         *time.Duration `json:"time_taken,omitempty"`
	AgentId           *string        `json:"agent_id,omitempty"`
	SessionId         string         `json:"session_id"`
}

type ChatSuggestion struct {
	ID         string `json:"id"`
	Suggestion string `json:"suggestion"`
}

type ChatClarification struct {
	ID        string  `json:"id"`
	Questions string  `json:"questions"`
	Answer    *string `json:"answer,omitempty"`
}
