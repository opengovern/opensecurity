package dag

import (
	"errors"
	"testing"
)

func TestTopSort(t *testing.T) {
	graph := initGraph()

	// a -> b -> c
	graph.AddEdge("a", "b")
	graph.AddEdge("b", "c")

	results, err := graph.TopologicalSort()
	if err != nil {
		t.Error(err)
		return
	}
	if results[0] != "c" || results[1] != "b" || results[2] != "a" {
		t.Errorf("Wrong sort order: %v", results)
	}
}

func TestTopSort2(t *testing.T) {
	graph := initGraph()

	// a -> c
	// a -> b
	// b -> c
	graph.AddEdge("a", "c")
	graph.AddEdge("a", "b")
	graph.AddEdge("b", "c")

	results, err := graph.TopologicalSort()
	if err != nil {
		t.Error(err)
		return
	}
	if results[0] != "c" || results[1] != "b" || results[2] != "a" {
		t.Errorf("Wrong sort order: %v", results)
	}
}

func TestTopSort3(t *testing.T) {
	graph := initGraph()

	// a -> b
	// a -> d
	// d -> c
	// c -> b
	graph.AddEdge("a", "b")
	graph.AddEdge("a", "d")
	graph.AddEdge("d", "c")
	graph.AddEdge("c", "b")

	results, err := graph.TopologicalSort()
	if err != nil {
		t.Error(err)
		return
	}
	if len(results) != 4 {
		t.Errorf("Wrong number of results: %v", results)
		return
	}
	expected := [4]string{"b", "c", "d", "a"}
	for i := 0; i < 4; i++ {
		if results[i] != expected[i] {
			t.Errorf("Wrong sort order: %v", results)
			break
		}
	}
}

func TestTopSortCycleError(t *testing.T) {
	graph := initGraph()

	// a -> b
	// b -> a
	graph.AddEdge("a", "b")
	graph.AddEdge("b", "a")

	_, err := graph.TopologicalSort()
	if err == nil {
		t.Errorf("Expected cycle error")
		return
	}
	if !errors.Is(err, ErrNotAcyclic) {
		t.Errorf("Error is not cycle error: %q", err)
	}
}

func TestTopSortCycleError2(t *testing.T) {
	graph := initGraph()

	// a -> b
	// b -> c
	// c -> a
	graph.AddEdge("a", "b")
	graph.AddEdge("b", "c")
	graph.AddEdge("c", "a")

	_, err := graph.TopologicalSort()
	if err == nil {
		t.Errorf("Expected cycle error")
		return
	}
	if !errors.Is(err, ErrNotAcyclic) {
		t.Errorf("Error is not cycle error: %q", err)
	}
}

func TestTopSortCycleError3(t *testing.T) {
	graph := initGraph()

	// a -> b
	// b -> c
	// c -> b
	graph.AddEdge("a", "b")
	graph.AddEdge("b", "c")
	graph.AddEdge("c", "b")

	_, err := graph.TopologicalSort()
	if err == nil {
		t.Errorf("Expected cycle error")
		return
	}
	if !errors.Is(err, ErrNotAcyclic) {
		t.Errorf("Error is not cycle error: %q", err)
	}
}

func initGraph() *DirectedAcyclicGraph {
	graph := NewDirectedAcyclicGraph()
	return graph
}
