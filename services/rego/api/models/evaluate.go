package models

import (
	"github.com/open-policy-agent/opa/rego"
)

type RegoEvaluateRequest struct {
	Policies []string `json:"policies"`
	Query    string   `json:"query"`
}

type RegoEvaluateResponse struct {
	Results rego.ResultSet `json:"result"`
}
