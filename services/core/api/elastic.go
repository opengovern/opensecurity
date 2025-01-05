package api

type BenchmarkTimeResponse struct {
	Aggregations BenchmarkTimeAggregation `json:"aggregations"`
}
type BenchmarkTimeAggregation struct {
	ReportTime ReportTimeAggregate `json:"reportTime"`
}
type ReportTimeAggregate struct {
	Buckets []Bucket `json:"buckets"`
}
type Bucket struct {
	Key      int64 `json:"key"`
	DocCount int64 `json:"doc_count"`
}
