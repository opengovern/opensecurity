package dag

import (
	"fmt"
)

var ErrNotAcyclic = fmt.Errorf("graph is not acyclic")

type DirectedAcyclicGraph struct {
	nodes map[string]node
}

type node struct {
	id    string
	edges map[string]struct{}
}

func (n node) addEdge(toNodeId string) {
	n.edges[toNodeId] = struct{}{}
}

func (n node) getEdgeTargets() []string {
	var nodeIds []string
	for e, _ := range n.edges {
		nodeIds = append(nodeIds, e)
	}
	return nodeIds
}

func newNode(id string) node {
	return node{
		id:    id,
		edges: make(map[string]struct{}),
	}
}

func NewDirectedAcyclicGraph() *DirectedAcyclicGraph {
	return &DirectedAcyclicGraph{
		nodes: make(map[string]node),
	}
}

func (g *DirectedAcyclicGraph) ContainsNode(nodeId string) bool {
	_, ok := g.nodes[nodeId]
	return ok
}

func (g *DirectedAcyclicGraph) AddNodeIdempotent(nodeId string) {
	if !g.ContainsNode(nodeId) {
		g.nodes[nodeId] = newNode(nodeId)
	}
}

func (g *DirectedAcyclicGraph) getOrAddNode(nodeId string) node {
	n, ok := g.nodes[nodeId]
	if !ok {
		n = newNode(nodeId)
		g.nodes[nodeId] = n
	}
	return n
}

func (g *DirectedAcyclicGraph) AddEdge(fromNodeId string, toNodeId string) {
	f := g.getOrAddNode(fromNodeId)
	g.AddNodeIdempotent(toNodeId)
	f.addEdge(toNodeId)
}

func (g *DirectedAcyclicGraph) TopologicalSort() ([]string, error) {
	rootNodes := make(map[string]struct{})
	for nodeId, _ := range g.nodes {
		rootNodes[nodeId] = struct{}{}
	}
	for _, n := range g.nodes {
		for _, edge := range n.getEdgeTargets() {
			delete(rootNodes, edge)
		}
	}
	if len(rootNodes) == 0 {
		return nil, ErrNotAcyclic
	}

	fmt.Println("rootNodes", rootNodes)

	result := newOrderedSet()
	visited := make(map[string]bool)
	for nodeId, _ := range rootNodes {
		if !visited[nodeId] {
			err := g.visit(nodeId, visited, result)
			if err != nil {
				return nil, err
			}
		}
	}
	return result.getElements(), nil
}

func (g *DirectedAcyclicGraph) visit(rootId string, visited map[string]bool, results *orderedSet) error {
	v := visited[rootId]
	if v {
		return ErrNotAcyclic
	}
	visited[rootId] = true

	n := g.nodes[rootId]

	for _, edge := range n.getEdgeTargets() {
		visitedClone := make(map[string]bool)
		for k, v := range visited {
			visitedClone[k] = v
		}
		err := g.visit(edge, visitedClone, results)
		if err != nil {
			return err
		}
	}

	results.add(rootId)
	return nil
}

type orderedSet struct {
	elements []string
	set      map[string]struct{}
}

func newOrderedSet() *orderedSet {
	return &orderedSet{
		elements: make([]string, 0),
		set:      make(map[string]struct{}),
	}
}

func (os *orderedSet) add(element string) {
	if _, ok := os.set[element]; !ok {
		os.elements = append(os.elements, element)
		os.set[element] = struct{}{}
	}
}

func (os *orderedSet) getElements() []string {
	return os.elements
}
