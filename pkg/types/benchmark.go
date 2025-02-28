package types

type FullBenchmark struct {
	ID    string `json:"ID" example:"azure_cis_v140"` // Benchmark ID
	Title string `json:"title" example:"CIS v1.4.0"`  // Benchmark title
}
