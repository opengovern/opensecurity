package shared

import "github.com/opengovern/opencomply/pkg/types"

type Query struct {
	QueryID        *string `json:"QueryID,omitempty" yaml:"QueryID,omitempty"`
	ID             string  `json:"ID,omitempty" yaml:"ID,omitempty"`
	Engine         string  `json:"Engine" yaml:"Engine"`
	QueryToExecute string  `json:"Definition" yaml:"Definition"`

	PrimaryTable *string          `json:"PrimaryResource" yaml:"PrimaryResource"`
	ListOfTables []string         `json:"ListOfResources" yaml:"ListOfResources"`
	Parameters   []QueryParameter `json:"Parameters" yaml:"Parameters"`
	Global       bool             `json:"Global,omitempty" yaml:"Global,omitempty"`

	RegoPolicies []string `json:"RegoPolicies,omitempty" yaml:"RegoPolicies,omitempty"`
}

type Policy struct {
	ID              *string              `json:"id,omitempty" yaml:"id,omitempty"`
	Title           string               `json:"title,omitempty" yaml:"title,omitempty"`
	Description     string               `json:"description,omitempty" yaml:"description,omitempty"`
	Ref             *string              `json:"@ref,omitempty" yaml:"@ref,omitempty"`
	Language        types.PolicyLanguage `json:"language" yaml:"language"`
	PrimaryResource string               `json:"primary_resource" yaml:"primary_resource"`
	ExampleData     *string              `json:"example_data,omitempty" yaml:"example_data,omitempty"`
	Definition      string               `json:"definition" yaml:"definition"`
	RegoPolicies    []string             `json:"rego_policies,omitempty" yaml:"RegoPolicies,omitempty"`
}

type ControlParameter struct {
	Key   string `json:"key" yaml:"key"`
	Value string `json:"value" yaml:"value"`
}

type QueryParameter struct {
	Key          string `json:"Key" yaml:"Key"`
	Required     bool   `json:"Required" yaml:"Required"`
	DefaultValue string `json:"DefaultValue" yaml:"DefaultValue"`
}

type ParameterDefaultValue struct {
	Key      string   `json:"key" yaml:"key"`
	Value    string   `json:"value" yaml:"value"`
	Controls []string `json:"controls" yaml:"controls"`
}

type ParameterDefaultValueFile struct {
	Type       string                  `json:"type" yaml:"type"`
	Parameters []ParameterDefaultValue `json:"parameters" yaml:"parameters"`
}
