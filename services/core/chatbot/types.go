package chatbot

// ResultType indicates the type of response from the LLM inference.
type ResultType string

const (
	ResultTypeSuccess             ResultType = "SUCCESS"
	ResultTypeClarificationNeeded ResultType = "CLARIFICATION_NEEDED"
	ResultTypeError               ResultType = "ERROR"
	ResultTypeMalformed           ResultType = "MALFORMED_RESPONSE" // For when parsing fails
	ResultTypeRawSQL              ResultType = "RAW_SQL"            // For fallback when JSON fails but SQL is found
)

// InferenceResult holds the structured result parsed from the LLM response.
type InferenceResult struct {
	Type ResultType `json:"type"` // Indicate the kind of result (not from JSON)

	// Fields for SUCCESS
	Query                     string   `json:"query,omitempty"`
	PrimaryInterpretation     string   `json:"primary_interpretation,omitempty"`
	AdditionalInterpretations []string `json:"additional_interpretations,omitempty"`

	// Fields for CLARIFICATION_NEEDED
	ClarifyingQuestions []string `json:"clarifying_questions,omitempty"`

	// Fields for ERROR
	Reason string `json:"reason,omitempty"` // Used if result is ERROR

	// Field for raw parsing issues or raw SQL fallback
	RawResponse string `json:"raw_response,omitempty"` // Store original if parsing fails or it's raw SQL
}

// --- Helper structs for intermediate JSON parsing ---

// BaseResponse is used to determine the 'result' type first.
type BaseResponse struct {
	Result ResultType `json:"result"` // Use ResultType for direct mapping
}

// SuccessResponse maps the fields for a SUCCESS result.
type SuccessResponse struct {
	Result                    ResultType `json:"result"`
	PrimaryInterpretation     string     `json:"primary_interpretation"`
	Query                     string     `json:"query"`
	AdditionalInterpretations []string   `json:"additional_interpretations"`
}

// ClarificationResponse maps the fields for a CLARIFICATION_NEEDED result.
type ClarificationResponse struct {
	Result              ResultType `json:"result"`
	ClarifyingQuestions []string   `json:"clarifying_questions"`
}

// ErrorResponse maps the fields for an ERROR result.
type ErrorResponse struct {
	Result ResultType `json:"result"`
	Reason string     `json:"reason"`
}
