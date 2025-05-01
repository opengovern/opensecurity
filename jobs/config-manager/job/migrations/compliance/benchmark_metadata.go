package compliance

import (
	"github.com/opengovern/opensecurity/services/compliance/db"
)

type BenchmarkControlsCache struct {
	Controls map[string]bool
}

// getControlsUnderBenchmark ctx context.Context, benchmarkId string -> primaryTables, listOfTables, error
func getControlsUnderBenchmark(benchmark db.Benchmark, benchmarksCache map[string]BenchmarkControlsCache) (map[string]bool, error) {
	controls := make(map[string]bool)

	var err error
	for _, c := range benchmark.Controls {
		controls[c.ID] = true
	}

	for _, child := range benchmark.Children {
		var childControls map[string]bool
		if cache, ok := benchmarksCache[child.ID]; ok {
			childControls = cache.Controls
		} else {
			childControls, err = getControlsUnderBenchmark(child, benchmarksCache)
			if err != nil {
				return nil, err
			}
			benchmarksCache[child.ID] = BenchmarkControlsCache{Controls: childControls}
		}
		for k, _ := range childControls {
			controls[k] = true
		}
	}
	return controls, nil
}

type BenchmarkTablesCache struct {
	PrimaryTables map[string]bool
	ListTables    map[string]bool
}

// getTablesUnderBenchmark ctx context.Context, benchmarkId string -> primaryTables, listOfTables, error
func (g *GitParser) getTablesUnderBenchmark(benchmark db.Benchmark, benchmarkCache map[string]BenchmarkTablesCache) (map[string]bool, map[string]bool, error) {
	primaryTables := make(map[string]bool)
	listOfTables := make(map[string]bool)

	for _, c := range benchmark.Controls {
		if query, ok := g.controlsPolicies[c.ID]; ok {
			if query.PrimaryResource != "" {
				primaryTables[query.PrimaryResource] = true
			}
			for _, t := range query.ListOfResources {
				if t == "" {
					continue
				}
				listOfTables[t] = true
			}
		}
	}

	var err error
	for _, child := range benchmark.Children {
		var childPrimaryTables, childListOfTables map[string]bool
		if cache, ok := benchmarkCache[child.ID]; ok {
			childPrimaryTables = cache.PrimaryTables
			childListOfTables = cache.ListTables
		} else {
			childPrimaryTables, childListOfTables, err = g.getTablesUnderBenchmark(child, benchmarkCache)
			if err != nil {
				return nil, nil, err
			}
			benchmarkCache[child.ID] = BenchmarkTablesCache{
				PrimaryTables: childPrimaryTables,
				ListTables:    childListOfTables,
			}
		}

		for k, _ := range childPrimaryTables {
			primaryTables[k] = true
		}
		for k, _ := range childListOfTables {
			childListOfTables[k] = true
		}
	}
	return primaryTables, listOfTables, nil
}
