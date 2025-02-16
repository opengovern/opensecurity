package models

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/lib/pq"
	"github.com/opengovern/og-util/pkg/integration"
	api "github.com/opengovern/opencomply/services/integration/api/models"
)

type Integration struct {
	integration.Integration
}

func (i *Integration) AddLabel(key, value string) (*pgtype.JSONB, error) {
	var labels map[string]string
	if i.Labels.Status == pgtype.Present {
		if err := json.Unmarshal(i.Labels.Bytes, &labels); err != nil {
			return nil, err
		}
	} else {
		labels = make(map[string]string)
	}

	labels[key] = value

	labelsJsonData, err := json.Marshal(labels)
	if err != nil {
		return nil, err
	}
	integrationLabelsJsonb := pgtype.JSONB{}
	err = integrationLabelsJsonb.Set(labelsJsonData)
	i.Labels = integrationLabelsJsonb

	return &integrationLabelsJsonb, nil
}

func (i *Integration) AddAnnotations(key, value string) (*pgtype.JSONB, error) {
	var annotation map[string]string
	if i.Annotations.Status == pgtype.Present {
		if err := json.Unmarshal(i.Annotations.Bytes, &annotation); err != nil {
			return nil, err
		}
	} else {
		annotation = make(map[string]string)
	}

	annotation[key] = value

	annotationsJsonData, err := json.Marshal(annotation)
	if err != nil {
		return nil, err
	}
	integrationAnnotationsJsonb := pgtype.JSONB{}
	err = integrationAnnotationsJsonb.Set(annotationsJsonData)
	i.Annotations = integrationAnnotationsJsonb

	return &integrationAnnotationsJsonb, nil
}

func (i *Integration) ToApi() (*api.Integration, error) {
	var labels map[string]string
	if i.Labels.Status == pgtype.Present {
		if err := json.Unmarshal(i.Labels.Bytes, &labels); err != nil {
			return nil, err
		}
	}

	var annotations map[string]string
	if i.Annotations.Status == pgtype.Present {
		if err := json.Unmarshal(i.Annotations.Bytes, &annotations); err != nil {
			return nil, err
		}
	}

	return &api.Integration{
		IntegrationID:   i.IntegrationID.String(),
		Name:            i.Name,
		ProviderID:      i.ProviderID,
		IntegrationType: i.IntegrationType,
		CredentialID:    i.CredentialID.String(),
		State:           api.IntegrationState(i.State),
		LastCheck:       i.LastCheck,
		Labels:          labels,
		Annotations:     annotations,
	}, nil
}

type IntegrationResourcetypes struct {
	IntegrationID uuid.UUID      `gorm:"primaryKey;type:uuid"` // Auto-generated UUID
	ResourceTypes pq.StringArray `gorm:"type:text[]"`
}
