package models

func (qp PolicyParameterValues) GetKey() string {
	return qp.Key
}

func (qp PolicyParameterValues) GetValue() string {
	return qp.Value
}
