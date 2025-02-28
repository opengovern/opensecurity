package types

type PolicyLanguage string

const (
	PolicyLanguageSQL       PolicyLanguage = "sql"
	PolicyLanguageRego      PolicyLanguage = "rego"
	PolicyLanguageUndefined PolicyLanguage = ""
)
