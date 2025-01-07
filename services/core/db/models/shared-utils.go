package models

import "github.com/opengovern/opencomply/services/core/api"

func (qp PolicyParameterValues) GetKey() string {
	return qp.Key
}

func (qp PolicyParameterValues) GetValue() string {
	return qp.Value
}

func (qp PolicyParameterValues) ToAPI() api.QueryParameter {
	return api.QueryParameter{
		Key:       qp.Key,
		ControlID: qp.ControlID,
		Value:     qp.Value,
	}
}

func QueryParameterFromAPI(apiQP api.QueryParameter) PolicyParameterValues {
	var qp PolicyParameterValues
	qp.Key = apiQP.Key
	qp.Value = apiQP.Value
	qp.ControlID = apiQP.ControlID
	return qp
}

func (qp QueryParameter) ToApi() api.QueryParameter {
	return api.QueryParameter{
		Key:      qp.Key,
		Required: qp.Required,
	}
}

func (q Query) ToApi() api.Query {
	query := api.Query{
		ID:             q.ID,
		QueryToExecute: q.QueryToExecute,
		ListOfTables:   q.ListOfTables,
		PrimaryTable:   q.PrimaryTable,
		Engine:         q.Engine,
		Parameters:     make([]api.QueryParameter, 0, len(q.Parameters)),
		Global:         q.Global,
		CreatedAt:      q.CreatedAt,
		UpdatedAt:      q.UpdatedAt,
	}
	for _, p := range q.Parameters {
		query.Parameters = append(query.Parameters, p.ToApi())
	}
	return query
}
