package query_validator

type QueryResult struct {
	Headers []string `json:"headers"` // Column names
	Result  [][]any  `json:"result"`  // Result of query. in order to access a specific cell please use Result[Row][Column]
}
