package describe

import (
	"github.com/opengovern/og-util/pkg/source"
)

type SourceAction string

const (
	SourceDelete SourceAction = "DELETE"
)

type SourceEvent struct {
	Action     SourceAction
	SourceID   string
	AccountID  string
	SourceType source.Type
	Secret     string
}
