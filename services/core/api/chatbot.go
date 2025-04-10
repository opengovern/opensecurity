package api

// QueryAttempt represents an attempt to generate and validate an SQL query.
type QueryAttempt struct {
	Query string `json:"query"`
	Error string `json:"error"`
}
type GenerateQueryRequest struct {
	Question         string         `json:"question"`
	PreviousAttempts []QueryAttempt `json:"previous_attempts"`
	Agent            *string        `json:"agent,omitempty"`
}

type GenerateQueryResponse struct {
	Query string `json:"query"`
	Agent string `json:"agent"`
}
