package compliance

import (
	"github.com/opengovern/og-util/pkg/integration"
	"github.com/opengovern/opencomply/jobs/post-install-job/job/migrations/shared"
)

type FrameworkFile struct {
	Framework Framework `json:"framework" yaml:"framework"`
}

type ControlGroupFile struct {
	ControlGroup Framework `json:"control-group" yaml:"control-group"`
}

type FrameworkMetadata struct {
	Defaults struct {
		AutoAssign        *bool `json:"auto-assign" yaml:"auto-assign"`
		Enabled           bool  `json:"enabled" yaml:"enabled"`
		TracksDriftEvents bool  `json:"tracks-drift-events" yaml:"tracks-drift-events"`
	} `json:"defaults"`
	Tags map[string][]string `json:"tags" yaml:"tags"`
}

type Framework struct {
	ID           string              `json:"id" yaml:"id"`
	Title        string              `json:"title" yaml:"title"`
	Description  string              `json:"description" yaml:"description"`
	SectionCode  string              `json:"section-code" yaml:"section-code"`
	Metadata     *FrameworkMetadata  `json:"metadata" yaml:"metadata"`
	Tags         map[string][]string `json:"tags" yaml:"tags"`
	ControlGroup []Framework         `json:"control-group" yaml:"control-group"`
	Controls     []string            `json:"controls" yaml:"controls"`
}

type Benchmark struct {
	ID                string              `json:"ID" yaml:"ID"`
	Title             string              `json:"Title" yaml:"Title"`
	SectionCode       string              `json:"SectionCode" yaml:"SectionCode"`
	Description       string              `json:"Description" yaml:"Description"`
	Children          []string            `json:"Children" yaml:"Children"`
	Tags              map[string][]string `json:"Tags" yaml:"Tags"`
	AutoAssign        *bool               `json:"AutoAssign" yaml:"AutoAssign"`
	Enabled           bool                `json:"Enabled" yaml:"Enabled"`
	TracksDriftEvents bool                `json:"TracksDriftEvents" yaml:"TracksDriftEvents"`
	Controls          []string            `json:"Controls" yaml:"Controls"`
}

type Control struct {
	ID              string                    `json:"id" yaml:"id"`
	Title           string                    `json:"title" yaml:"title"`
	Description     string                    `json:"description" yaml:"description"`
	IntegrationType []string                  `json:"integration_type" yaml:"integration_type"`
	Parameters      []shared.ControlParameter `json:"parameters" yaml:"parameters"`
	Policy          *shared.Policy            `json:"policy" yaml:"policy"`
	Severity        string                    `json:"severity" yaml:"severity"`
	Tags            map[string][]string       `json:"tags" yaml:"tags"`
}

type QueryView struct {
	ID          string              `json:"id" yaml:"id"`
	Title       string              `json:"title" yaml:"title"`
	Description string              `json:"description" yaml:"description"`
	Query       string              `json:"query" yaml:"query"`
	Tags        map[string][]string `json:"tags" yaml:"tags"`
}

type NamedQuery struct {
	ID               string                    `json:"id" yaml:"id"`
	Title            string                    `json:"title" yaml:"title"`
	Description      string                    `json:"description" yaml:"description"`
	Parameters       []shared.ControlParameter `json:"parameters" yaml:"parameters"`
	IntegrationTypes []integration.Type        `json:"integration_type" yaml:"integration_type"`
	Query            string                    `json:"query" yaml:"query"`
	Tags             map[string][]string       `json:"tags" yaml:"tags"`
}
