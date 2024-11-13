package dag

import (
	"fmt"
	"sync"

	"github.com/aacanakin/dag/queue"
	"github.com/aacanakin/dag/set"
	"github.com/aacanakin/dag/stack"
	"github.com/pkg/errors"
)

// Vertex represents a vertex or node in the graph
type Vertex = string

// keys represent vertex & values are direct edges to vertices
// Edges represents the edges of the graph
type Edges map[Vertex][]Vertex

// New creates an empty graph with no vertices & edges and returns it
func New() *Graph {
	return &Graph{
		vertices: []Vertex{},
		edges:    map[Vertex][]Vertex{},
	}
}

// Graph represents a directed asyclic graph
type Graph struct {
	mu       sync.RWMutex
	vertices []Vertex
	edges    Edges
}

// Edges returns the edges of the graph
func (g *Graph) Edges() Edges {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.edges
}

// Vertices returns the vertices of the graph
func (g *Graph) Vertices() []Vertex {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.vertices
}

// Exists checks if a vertex exists in the graph
func (g *Graph) Exists(vertex Vertex) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	_, exists := g.edges[vertex]
	return exists
}

// Prev returns the previous vertices of a given vertex
func (g *Graph) Prev(vertex Vertex) (prev []Vertex, err error) {
	if existing := g.Exists(vertex); !existing {
		return []Vertex{}, fmt.Errorf("vertex %s is not found in graph", vertex)
	}

	prev = []Vertex{}

	for _, v := range g.Vertices() {
		nextVertices, err := g.Next(v)
		if err != nil {
			return []Vertex{}, errors.Wrap(err, "could not calculate prev")
		}
		if some(nextVertices, func(nextVertex Vertex) bool {
			return nextVertex == vertex
		}) {
			prev = append(prev, v)
		}
	}

	return prev, nil
}

// Next returns the next vertices of a given vertex
func (g *Graph) Next(vertex Vertex) ([]Vertex, error) {
	if existing := g.Exists(vertex); !existing {
		return []Vertex{}, fmt.Errorf("vertex %s is not found in graph", vertex)
	}

	g.mu.RLock()
	defer g.mu.RUnlock()
	next := g.edges[vertex]
	if next == nil {
		return []Vertex{}, nil
	}
	return next, nil
}

// ReverseEdges returns the reverse edges of the graph
// example: {a: [b, c], b: [c], c: []} -> {a: [], b: [a], c: [a, b]}
func (g *Graph) ReverseEdges() (Edges, error) {
	reverse := map[Vertex][]Vertex{}
	for _, vertex := range g.Vertices() {
		reverse[vertex] = []Vertex{}
	}

	for _, vertex := range g.Vertices() {
		nextVertices, err := g.Next(vertex)
		if err != nil {
			return nil, errors.Wrap(err, "could not reverse edges")
		}
		for _, nextVertex := range nextVertices {
			reverse[nextVertex] = append(reverse[nextVertex], vertex)
		}
	}

	return reverse, nil
}

// Reverse returns the reverse graph of the graph with reversed edges
func (g *Graph) Reverse() (*Graph, error) {
	revEdges, err := g.ReverseEdges()
	if err != nil {
		return nil, errors.Wrap(err, "could not reverse graph")
	}

	return &Graph{vertices: g.Vertices(), edges: revEdges}, nil
}

// DFS performs depth first search on the graph starting from the given vertex
func (g *Graph) DFS(start Vertex) (result []Vertex, err error) {
	stack := stack.New()
	stack.Push(start)

	result = []Vertex{}
	visited := set.New()

	for !stack.IsEmpty() {
		current, err := stack.Pop()
		if err != nil {
			return []Vertex{}, errors.Wrap(err, "could not perform dfs on the graph")
		}

		if !visited.Has(current) {
			visited.Add(current)
			result = append(result, current)

			next, err := g.Next(current)
			if err != nil {
				return []Vertex{}, errors.Wrap(err, "could not perform dfs on the graph")
			}
			for _, nextVertex := range next {
				stack.Push(nextVertex)
			}
		}
	}

	return result, nil
}

// BFS performs breadth first search on the graph starting from the given vertex
func (g *Graph) BFS(start Vertex) (result []Vertex, err error) {
	queue := queue.New()
	queue.Enqueue(start)
	visited := set.New()

	for queue.Size() > 0 {
		current, err := queue.Pop()
		if err != nil {
			return []Vertex{}, err
		}

		if !visited.Has(current) {
			visited.Add(current)
			result = append(result, current)
		}

		next, err := g.Next(current)
		if err != nil {
			return []Vertex{}, errors.Wrap(err, "could not perform bfs on the graph")
		}
		if len(next) != 0 {
			for _, n := range next {
				queue.Enqueue(n)
			}
		}
	}

	return result, nil
}

// Deps returns the dependencies of a vertex given vertex, in topological order
func (g *Graph) Deps(vertex string) (result []Vertex, err error) {
	reverse, err := g.Reverse()
	if err != nil {
		return []Vertex{}, errors.Wrap(err, fmt.Sprintf("could not calculate deps for vertex %s", vertex))
	}

	dfs, err := reverse.DFS(vertex)
	if err != nil {
		return []Vertex{}, errors.Wrap(err, fmt.Sprintf("could not calculate deps for vertex %s", vertex))
	}

	subgraph, err := g.SubGraph(exclude(dfs, vertex))
	if err != nil {
		return []Vertex{}, errors.Wrap(err, fmt.Sprintf("could not calculate deps for vertex %s", vertex))
	}

	sorted, err := subgraph.TopSort()
	if err != nil {
		return []Vertex{}, errors.Wrap(err, fmt.Sprintf("could not calculate deps for vertex %s", vertex))
	}

	return sorted, err
}

// ReverseDeps returns the reverse dependencies of a vertex given vertex, in topological order
func (g *Graph) ReverseDeps(vertex Vertex) (result []Vertex, err error) {
	dfs, err := g.DFS(vertex)
	if err != nil {
		return []Vertex{}, errors.Wrap(err, fmt.Sprintf("could not calculate reverse deps for vertex %s", vertex))
	}

	subgraph, err := g.SubGraph(exclude(dfs, vertex))
	if err != nil {
		return []Vertex{}, errors.Wrap(err, fmt.Sprintf("could not calculate reverse deps for vertex %s", vertex))
	}

	sorted, err := subgraph.TopSort()
	if err != nil {
		return []Vertex{}, errors.Wrap(err, fmt.Sprintf("could not calculate reverse deps for vertex %s", vertex))
	}

	return sorted, err
}

// Leaves returns the leaf vertices of the graph
func (g *Graph) Leaves() (leaves []Vertex, err error) {
	for _, vertex := range g.Vertices() {
		next, err := g.Next(vertex)
		if err != nil {
			return nil, errors.Wrap(err, "could not find leaves of graph")
		}
		if len(next) == 0 {
			leaves = append(leaves, vertex)
		}
	}

	return leaves, nil
}

// Roots returns the root vertices of the graph
func (g *Graph) Roots() ([]Vertex, error) {
	reverse, err := g.Reverse()
	if err != nil {
		return []Vertex{}, errors.Wrap(err, "could not find roots")
	}
	return reverse.Leaves()
}

// Append adds a new vertex to graph given vertex and previous vertices,
// returns error if any of the previous vertices is not present in graph
func (g *Graph) Append(v Vertex, prevVertices []Vertex) error {
	if existing := g.Exists(v); existing {
		return errors.Wrap(fmt.Errorf("duplicate node id=%s are not allowed", v), "could not append node to graph")
	}

	for _, prevVertex := range prevVertices {
		if !g.Exists(prevVertex) {
			return errors.Wrap(fmt.Errorf("prev vertex %s is not found in graph", prevVertex), "could not append node to graph")
		}
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	g.vertices = append(g.vertices, v)
	g.edges[v] = []Vertex{}

	for _, prevVertex := range prevVertices {
		g.edges[prevVertex] = append(g.edges[prevVertex], v)
	}

	return nil
}

// Add appends an unconnected node to the graph
func (g *Graph) Add(vertices ...Vertex) error {
	if len(vertices) == 0 {
		return fmt.Errorf("no vertices to add to graph")
	}
	for _, v := range vertices {
		if existing := g.Exists(v); existing {
			return fmt.Errorf("vertex %s already added. vertices must be unique", v)
		}

		g.mu.Lock()

		g.vertices = append(g.vertices, v)
		g.edges[v] = []Vertex{}

		g.mu.Unlock()
	}
	return nil
}

func (g *Graph) hasNext(from Vertex, to Vertex) (bool, error) {
	next, err := g.Next(from)
	if err != nil {
		return false, errors.Wrap(err, "could not perform hasNext")
	}
	return some(next, func(v Vertex) bool {
		return v == to
	}), nil
}

func (g *Graph) hasDep(from Vertex, to Vertex) bool {
	dfsVertices, err := g.DFS(from)
	if err != nil {
		return false
	}

	return some(dfsVertices, func(v Vertex) bool {
		return v == to
	})
}

// Connect connects two vertices in the graph
//
// returns error if;
//
// 	the edge already exists or if the edge creates a cycle
//
// 	the from vertex is not found in the graph
//
// 	the to vertex is not found in the graph
//
// 	the from vertex is the same as the to vertex
//
// it can be used to lazily initialize vertice connections
func (g *Graph) Connect(from Vertex, to Vertex) error {

	hasEdge, err := g.hasNext(from, to)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("could not connect vertex %s to vertex %s", from, to))
	}
	if hasEdge {
		return fmt.Errorf("could not connect vertex %s to vertex %s. edge already exists", from, to)
	}

	if hasCycle := g.hasDep(to, from); hasCycle {
		return fmt.Errorf("could not connect nodes. reason: cyclic edges are not allowed from %s to %s", from, to)
	}

	g.mu.Lock()
	defer g.mu.Unlock()
	g.edges[from] = append(g.edges[from], to)

	return nil
}

// DisconnectEdge disconnects two vertices in the graph
// returns error if the edge does not exist
func (g *Graph) DisconnectEdge(from Vertex, to Vertex) error {
	g.mu.RLock()
	edgeIndex := index(g.edges[from], func(v Vertex) bool {
		return v == to
	})
	g.mu.RUnlock()

	if edgeIndex < 0 {
		return fmt.Errorf("could not disconnect graph node prev=%s next=%s. edge does not exist", from, to)
	}

	g.mu.Lock()
	defer g.mu.Unlock()
	g.edges[from] = exclude(g.edges[from], to)

	return nil
}

// Disconnect disconnects all edges from and to a vertex
// returns error if the vertex is not found in the graph
func (g *Graph) Disconnect(v Vertex) error {
	prev, err := g.Prev(v)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("could not disconnect vertex %s", v))
	}

	// clear previous edges
	for _, prevVertex := range prev {
		if err := g.DisconnectEdge(prevVertex, v); err != nil {
			return err
		}
	}

	next, err := g.Next(v)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("could not disconnect vertex %s", v))
	}
	// clear next edges
	for _, nextVertex := range next {
		if err := g.DisconnectEdge(v, nextVertex); err != nil {
			return err
		}
	}

	return nil
}

// Remove removes node, all next nodes that are connected to that node & clears all edges that are related to node & deps
func (g *Graph) Remove(v Vertex) (removed []Vertex, err error) {
	toRemove, err := g.DFS(v)
	if err != nil {
		return removed, errors.Wrap(err, "could not remove node")
	}

	for _, removedVertex := range toRemove {
		err = g.Disconnect(removedVertex)
		if err != nil {
			return removed, errors.Wrap(err, "could not remove node")
		}
		removed = append(removed, removedVertex)
	}

	g.mu.Lock()
	defer g.mu.Unlock()
	g.vertices = exclude(g.vertices, toRemove...)
	delete(g.edges, v)

	return removed, nil
}

// TopSort applies topological sort algorithm to graph and returns vertices slice
func (g *Graph) TopSort() (result []Vertex, err error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	inDegree := make(map[Vertex]int, len(g.vertices))
	for _, vertex := range g.vertices {
		inDegree[vertex] = 0
	}

	for _, nextVertices := range g.edges {
		for _, nextVertex := range nextVertices {
			inDegree[nextVertex]++
		}
	}

	queue := queue.New()
	for _, vertex := range g.vertices {
		degree := inDegree[vertex]
		if degree == 0 {
			queue.Enqueue(vertex)
		}
	}

	for queue.Size() > 0 {

		vertex, err := queue.Pop()
		if err != nil {
			return []Vertex{}, errors.Wrap(err, "could not sort vertices")
		}

		result = append(result, vertex)

		for _, nextVertex := range g.edges[vertex] {
			inDegree[nextVertex]--

			if inDegree[nextVertex] == 0 {
				queue.Enqueue(nextVertex)
			}
		}
	}

	return result, nil
}

// DeepCopy creates a deep copy of the graph
func (g *Graph) DeepCopy() (*Graph, error) {
	graph := New()
	for _, vertex := range g.Vertices() {
		err := graph.Add(vertex)
		if err != nil {
			return nil, errors.Wrap(err, "could not create deep copy")
		}
	}

	for vertex, nextVertices := range g.Edges() {
		for _, nextVertex := range nextVertices {
			err := graph.Connect(vertex, nextVertex)
			if err != nil {
				return nil, errors.Wrap(err, "could not create deep copy")
			}
		}
	}

	return graph, nil
}

// SubGraph returns a subgraph of existing graph that includes input vertices, clips out vertices & non connected edges
func (g *Graph) SubGraph(vertices []Vertex) (graph *Graph, err error) {
	excludedVertices := exclude(g.Vertices(), vertices...)

	subGraph, err := g.DeepCopy()
	if err != nil {
		return nil, errors.Wrap(err, "could not create sub graph")
	}
	for _, vertex := range excludedVertices {
		err = subGraph.Disconnect(vertex)
		if err != nil {
			return nil, errors.Wrap(err, "could not create sub graph")
		}

		_, err = subGraph.Remove(vertex)
		if err != nil {
			return nil, errors.Wrap(err, "could not create sub graph")
		}
	}

	return subGraph, nil
}

func exclude(vertices []Vertex, exclude ...Vertex) []Vertex {
	return filter(vertices, func(vertice Vertex) bool {
		return !includes(exclude, vertice)
	})
}
