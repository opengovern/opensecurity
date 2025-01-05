package models

import "github.com/opengovern/opencomply/services/metadata/api"

type PolicyParameterValues struct {
	Key       string `gorm:"primaryKey"`
	ControlID string `gorm:"primaryKey"`
	Value     string `gorm:"type:text;not null"`
}

func (qp PolicyParameterValues) GetKey() string {
	return qp.Key
}

func (qp PolicyParameterValues) GetValue() string {
	return qp.Value
}

func (qp PolicyParameterValues) ToAPI() api.QueryParameter {
	return api.QueryParameter{
		Key:   qp.Key,
		Value: qp.Value,
	}
}

func QueryParameterFromAPI(apiQP api.QueryParameter) PolicyParameterValues {
	var qp PolicyParameterValues
	qp.Key = apiQP.Key
	qp.Value = apiQP.Value
	return qp
}
