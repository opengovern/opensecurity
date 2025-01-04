package models

import "github.com/opengovern/opencomply/services/metadata/api"





func (qp QueryParameterValues) GetKey() string {
	return qp.Key
}

func (qp QueryParameterValues) GetValue() string {
	return qp.Value
}

func (qp QueryParameterValues) ToAPI() api.QueryParameter {
	return api.QueryParameter{
		Key:   qp.Key,
		Value: qp.Value,
	}
}

func QueryParameterFromAPI(apiQP api.QueryParameter) QueryParameterValues {
	var qp QueryParameterValues
	qp.Key = apiQP.Key
	qp.Value = apiQP.Value
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