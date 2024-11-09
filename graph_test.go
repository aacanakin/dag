package dag_test

import (
	"fmt"
	"testing"

	"github.com/aacanakin/dag"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

/*
A -- B -- C
|    |
D -- E -- F
*/
func createGraph() *dag.Graph {
	g := dag.New()

	err := g.Add("A", "B", "C", "D", "E", "F")

	if err != nil {
		panic(errors.Wrap(err, "could not create sample graph for testing"))
	}

	edges := []struct {
		from string
		to   string
	}{
		{from: "A", to: "B"},
		{from: "A", to: "D"},
		{from: "B", to: "C"},
		{from: "B", to: "E"},
		{from: "D", to: "E"},
		{from: "E", to: "F"},
	}

	for _, edge := range edges {
		err := g.Connect(edge.from, edge.to)
		if err != nil {
			panic(errors.Wrap(err, fmt.Sprintf("could not connect %s to %s", edge.from, edge.to)))
		}
	}

	return g
}

func TestGraph(t *testing.T) {
	t.Run("Vertices", func(t *testing.T) {
		t.Run("should return all vertices of graph", func(t *testing.T) {
			g := createGraph()
			vertices := g.Vertices()

			assert.Equal(t, 6, len(vertices))

			assert.Contains(t, vertices, dag.Vertex("A"))
			assert.Contains(t, vertices, dag.Vertex("B"))
			assert.Contains(t, vertices, dag.Vertex("C"))
			assert.Contains(t, vertices, dag.Vertex("D"))
			assert.Contains(t, vertices, dag.Vertex("E"))
			assert.Contains(t, vertices, dag.Vertex("F"))
		})

		t.Run("should return empty slice for graph with no vertices", func(t *testing.T) {
			assert.Equal(t, 0, len(dag.New().Vertices()))
			assert.Equal(t, []dag.Vertex{}, dag.New().Vertices())
		})

		t.Run("should return true for existing vertex", func(t *testing.T) {
			g := createGraph()
			assert.Equal(t, true, g.Exists("D"))
		})
	})

	t.Run("Exists", func(t *testing.T) {
		t.Run("Exists", func(t *testing.T) {
			t.Run("should return false for non existing vertex", func(t *testing.T) {
				g := createGraph()
				assert.Equal(t, false, g.Exists("X"))
			})
		})
	})

	t.Run("Prev", func(t *testing.T) {

		type prevResult struct {
			prev []dag.Vertex
			err  error
		}

		g := createGraph()

		expected := map[string]prevResult{
			"A":       {prev: []dag.Vertex{}, err: nil},
			"B":       {prev: []dag.Vertex{"A"}, err: nil},
			"C":       {prev: []dag.Vertex{"B"}, err: nil},
			"D":       {prev: []dag.Vertex{"A"}, err: nil},
			"E":       {prev: []dag.Vertex{"B", "D"}, err: nil},
			"F":       {prev: []dag.Vertex{"E"}, err: nil},
			"invalid": {prev: []dag.Vertex{}, err: fmt.Errorf("vertex invalid is not found in graph")},
		}

		results := map[string]prevResult{}
		for vertex := range expected {
			prev, err := g.Prev(vertex)
			results[vertex] = prevResult{
				prev: prev,
				err:  err,
			}
		}

		for vertex, result := range expected {
			assert.Equal(t, result.prev, results[vertex].prev, fmt.Sprintf("Checking prev of vertex %s", vertex))
			assert.Equal(t, result.err, results[vertex].err, fmt.Sprintf("Checking error of vertex %s", vertex))
		}

	})

	t.Run("Next", func(t *testing.T) {

		type nextResult struct {
			next []dag.Vertex
			err  error
		}

		g := createGraph()

		expected := map[string]nextResult{
			"A":       {next: []dag.Vertex{"B", "D"}, err: nil},
			"B":       {next: []dag.Vertex{"C", "E"}, err: nil},
			"C":       {next: []dag.Vertex{}, err: nil},
			"D":       {next: []dag.Vertex{"E"}, err: nil},
			"E":       {next: []dag.Vertex{"F"}, err: nil},
			"F":       {next: []dag.Vertex{}, err: nil},
			"invalid": {next: []dag.Vertex{}, err: fmt.Errorf("vertex invalid is not found in graph")},
		}

		results := map[string]nextResult{}
		for vertex := range expected {
			next, err := g.Next(vertex)
			results[vertex] = nextResult{
				next: next,
				err:  err,
			}
		}

		for vertex, result := range expected {
			assert.Equal(t, result.next, results[vertex].next, fmt.Sprintf("Checking next of vertex %s", vertex))
			assert.Equal(t, result.err, results[vertex].err, fmt.Sprintf("Checking error of vertex %s", vertex))
		}
	})

	t.Run("ReverseEdges", func(t *testing.T) {
		t.Run("should return no edges with empty edges", func(t *testing.T) {
			g := dag.New()
			err := g.Add(dag.Vertex("A"))
			assert.Nil(t, err)

			revEdges, err := g.ReverseEdges()
			expected := dag.Edges{"A": {}}
			assert.Nil(t, err)
			assert.Equal(t, expected, revEdges)
		})

		t.Run("should reverse single edge correctly", func(t *testing.T) {
			g := dag.New()
			err := g.Add("A")
			assert.Nil(t, err)
			err = g.Add("B")
			assert.Nil(t, err)
			err = g.Connect("A", "B")
			assert.Nil(t, err)

			revEdges, err := g.ReverseEdges()
			expected := dag.Edges{"A": {}, "B": {"A"}}
			assert.Nil(t, err)
			assert.Equal(t, expected, revEdges)
		})

		t.Run("should reverse graph with multiple next for a single vertex correctly", func(t *testing.T) {
			g := dag.New()
			var err error
			err = g.Add(dag.Vertex("A"))
			assert.Nil(t, err)
			err = g.Add(dag.Vertex("B"))
			assert.Nil(t, err)
			err = g.Add(dag.Vertex("C"))
			assert.Nil(t, err)
			err = g.Connect("A", "B")
			assert.Nil(t, err)
			err = g.Connect("A", "C")
			assert.Nil(t, err)

			revEdges, err := g.ReverseEdges()
			expected := dag.Edges{"A": {}, "B": {"A"}, "C": {"A"}}
			assert.Nil(t, err)
			assert.Equal(t, expected, revEdges)
		})

		t.Run("should reverse sample graph", func(t *testing.T) {
			g := createGraph()
			revEdges, err := g.ReverseEdges()

			expected := dag.Edges(
				map[dag.Vertex][]dag.Vertex{
					"A": {},
					"B": {"A"},
					"C": {"B"},
					"D": {"A"},
					"E": {"B", "D"},
					"F": {"E"},
				},
			)

			assert.Nil(t, err)
			for v := range expected {
				assert.Equal(t, expected[v], revEdges[v], fmt.Sprintf("Checking vertex %s", v))
			}
		})
	})

	t.Run("Reverse", func(t *testing.T) {
		t.Run("should reverse sample graph", func(t *testing.T) {
			g := createGraph()
			reverse, err := g.Reverse()

			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"A", "B", "C", "D", "E", "F"}, reverse.Vertices())

			expectedReverseEdges := dag.Edges(
				map[dag.Vertex][]dag.Vertex{
					"A": {},
					"B": {"A"},
					"C": {"B"},
					"D": {"A"},
					"E": {"B", "D"},
					"F": {"E"},
				},
			)

			for v := range expectedReverseEdges {
				assert.Equal(t, expectedReverseEdges[v], reverse.Edges()[v], fmt.Sprintf("Checking vertex %s", v))
			}
		})
	})

	t.Run("DFS", func(t *testing.T) {
		t.Run("should return dfs for root vertex", func(t *testing.T) {
			g := createGraph()

			dfs, err := g.DFS("A")
			expected := []dag.Vertex{"A", "D", "E", "F", "B", "C"}

			assert.Nil(t, err)
			assert.Equal(t, expected, dfs)
		})

		t.Run("should return dfs for leaf vertex", func(t *testing.T) {
			g := createGraph()

			dfs, err := g.DFS("C")

			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"C"}, dfs)
		})

		t.Run("should return dfs for middle vertex", func(t *testing.T) {
			g := createGraph()

			dfs, err := g.DFS("B")

			assert.Nil(t, err)

			expected := []dag.Vertex{"B", "E", "F", "C"}
			assert.Equal(t, expected, dfs)
		})

		t.Run("should return error for non existing vertex", func(t *testing.T) {
			g := createGraph()

			dfs, err := g.DFS("X")

			assert.NotNil(t, err)
			assert.Equal(t, []dag.Vertex{}, dfs)
		})
	})

	t.Run("BFS", func(t *testing.T) {
		t.Run("should return bfs for root vertex", func(t *testing.T) {
			g := createGraph()

			bfs, err := g.BFS("A")

			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"A", "B", "D", "C", "E", "F"}, bfs)
		})

		t.Run("should return bfs for leaf vertex", func(t *testing.T) {
			g := createGraph()

			bfs, err := g.BFS("F")

			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"F"}, bfs)
		})

		t.Run("should return bfs for middle vertex", func(t *testing.T) {
			g := createGraph()

			bfs, err := g.BFS("B")

			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"B", "C", "E", "F"}, bfs)
		})

		t.Run("should return error for non existing vertex", func(t *testing.T) {
			g := createGraph()

			bfs, err := g.BFS("X")

			assert.NotNil(t, err)
			assert.Equal(t, []dag.Vertex{}, bfs)
		})
	})

	type depsResult struct {
		deps []dag.Vertex
		err  error
	}

	t.Run("Deps", func(t *testing.T) {

		t.Run("should return deps for every vertex in sample graph", func(t *testing.T) {
			g := createGraph()

			expected := map[string]depsResult{
				"A": {[]dag.Vertex(nil), nil},
				"B": {[]dag.Vertex{"A"}, nil},
				"C": {[]dag.Vertex{"A", "B"}, nil},
				"D": {[]dag.Vertex{"A"}, nil},
				"E": {[]dag.Vertex{"A", "B", "D"}, nil},
				"F": {[]dag.Vertex{"A", "B", "D", "E"}, nil},
			}

			results := map[string]depsResult{}
			for _, v := range g.Vertices() {
				deps, err := g.Deps(v)

				results[v] = depsResult{
					deps,
					err,
				}
			}

			for vertex, result := range expected {
				assert.Equal(t, result.deps, results[vertex].deps, fmt.Sprintf("Checking deps of vertex %s", vertex))
				assert.Equal(t, result.err, results[vertex].err, fmt.Sprintf("Checking error of vertex %s", vertex))
			}
		})
	})

	t.Run("ReverseDeps", func(t *testing.T) {
		t.Run("should return reverse deps for every vertex in sample graph", func(t *testing.T) {
			g := createGraph()

			results := map[string]depsResult{}
			expected := map[string]depsResult{
				"A": {[]dag.Vertex{"B", "D", "C", "E", "F"}, nil},
				"B": {[]dag.Vertex{"C", "E", "F"}, nil},
				"C": {[]dag.Vertex(nil), nil},
				"D": {[]dag.Vertex{"E", "F"}, nil},
				"E": {[]dag.Vertex{"F"}, nil},
				"F": {[]dag.Vertex(nil), nil},
			}

			for _, v := range g.Vertices() {
				deps, err := g.ReverseDeps(v)

				results[v] = depsResult{
					deps,
					err,
				}
			}

			for vertex, result := range expected {
				assert.Equal(t, result.deps, results[vertex].deps, fmt.Sprintf("Checking deps of vertex %s", vertex))
				assert.Equal(t, result.err, results[vertex].err, fmt.Sprintf("Checking error of vertex %s", vertex))
			}
		})
	})

	t.Run("Leaves", func(t *testing.T) {
		t.Run("should return leaf vertices for sample graph", func(t *testing.T) {
			g := createGraph()

			leaves, err := g.Leaves()

			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"C", "F"}, leaves)
		})

		t.Run("should return leaf vertices for reversed sample graph", func(t *testing.T) {
			g, err := createGraph().Reverse()

			assert.Nil(t, err)

			leaves, err := g.Leaves()

			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"A"}, leaves)
		})
	})

	t.Run("Roots", func(t *testing.T) {
		t.Run("should return root vertices for sample graph", func(t *testing.T) {
			g := createGraph()

			roots, err := g.Roots()

			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"A"}, roots)
		})

		t.Run("should return root vertices for reversed sample graph", func(t *testing.T) {
			g, err := createGraph().Reverse()

			assert.Nil(t, err)

			roots, err := g.Roots()

			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"C", "F"}, roots)
		})

		t.Run("should return multi disconnected root vertices", func(t *testing.T) {
			var err error
			g := dag.New()
			err = g.Add("A")
			assert.Nil(t, err)
			err = g.Add("B")
			assert.Nil(t, err)

			roots, err := g.Roots()
			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"A", "B"}, roots)
		})

		t.Run("should return multi interconnected root vertices", func(t *testing.T) {
			var err error
			g := dag.New()
			err = g.Add("A")
			assert.Nil(t, err)
			err = g.Add("B")
			assert.Nil(t, err)
			err = g.Add("C")
			assert.Nil(t, err)
			err = g.Add("D")
			assert.Nil(t, err)

			err = g.Connect("A", "B")
			assert.Nil(t, err)
			err = g.Connect("C", "D")
			assert.Nil(t, err)

			roots, err := g.Roots()
			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"A", "C"}, roots)
		})

		t.Run("should return multi connected root vertices", func(t *testing.T) {
			var err error
			g := dag.New()
			err = g.Add("R1")
			assert.Nil(t, err)
			err = g.Add("R2")
			assert.Nil(t, err)
			err = g.Add("A")
			assert.Nil(t, err)
			err = g.Add("B")
			assert.Nil(t, err)

			err = g.Connect("R1", "A")
			assert.Nil(t, err)
			err = g.Connect("R2", "A")
			assert.Nil(t, err)
			err = g.Connect("A", "B")
			assert.Nil(t, err)

			roots, err := g.Roots()
			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"R1", "R2"}, roots)
		})
	})

	t.Run("Append", func(t *testing.T) {
		t.Run("should append new vertex to sample graph", func(t *testing.T) {
			g := createGraph()

			err := g.Append("X", []dag.Vertex{"A", "B"})
			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"A", "B", "C", "D", "E", "F", "X"}, g.Vertices(), "Checking vertices after appending X")

			prev, err := g.Prev("X")
			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"A", "B"}, prev, fmt.Sprintf("Checking prev of X: %v", prev))

			next, err := g.Next("X")
			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{}, next, fmt.Sprintf("Checking next of X: %v", next))

			nextOfA, err := g.Next("A")
			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"B", "D", "X"}, nextOfA, fmt.Sprintf("Checking next of A: %v", nextOfA))

			nextOfB, err := g.Next("B")
			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"C", "E", "X"}, nextOfB, fmt.Sprintf("Checking next of B: %v", nextOfB))
		})

		t.Run("should append new vertex without any prev vertice", func(t *testing.T) {
			g := createGraph()

			err := g.Append("X", []dag.Vertex{})
			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"A", "B", "C", "D", "E", "F", "X"}, g.Vertices(), "Checking vertices")

			prev, err := g.Prev("X")
			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{}, prev, fmt.Sprintf("Checking prev of X: %v", prev))

			next, err := g.Next("X")
			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{}, next, fmt.Sprintf("Checking next of X: %v", next))

			nextOfA, err := g.Next("A")
			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"B", "D"}, nextOfA, fmt.Sprintf("Checking next of A: %v", nextOfA))

			nextOfB, err := g.Next("B")
			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"C", "E"}, nextOfB, fmt.Sprintf("Checking next of B: %v", nextOfB))
		})

		t.Run("should append new vertex with every vertex as prev", func(t *testing.T) {
			g := createGraph()

			err := g.Append("X", []dag.Vertex{"A", "B", "C", "D", "E", "F"})
			assert.Nil(t, err)

			assert.Equal(t, []dag.Vertex{"A", "B", "C", "D", "E", "F", "X"}, g.Vertices(), "Checking vertices")

			prev, err := g.Prev("X")
			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"A", "B", "C", "D", "E", "F"}, prev, fmt.Sprintf("Checking prev of X: %v", prev))

			next, err := g.Next("X")
			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{}, next, fmt.Sprintf("Checking next of X: %v", next))
		})

		t.Run("should return error for existing vertex", func(t *testing.T) {
			g := createGraph()

			err := g.Append("A", []dag.Vertex{})

			assert.NotNil(t, err)
		})

		t.Run("should return error for non existing prev vertices", func(t *testing.T) {
			g := createGraph()

			err := g.Append("X", []dag.Vertex{"Y"})

			assert.NotNil(t, err)
		})
	})

	t.Run("Add", func(t *testing.T) {
		t.Run("should add new vertex to sample graph", func(t *testing.T) {
			g := createGraph()

			err := g.Add("X")

			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"A", "B", "C", "D", "E", "F", "X"}, g.Vertices(), "Checking vertices")
		})

		t.Run("should return error while adding existing vertex to sample graph", func(t *testing.T) {
			g := createGraph()

			err := g.Add("A")

			assert.NotNil(t, err)
			assert.Equal(t, []dag.Vertex{"A", "B", "C", "D", "E", "F"}, g.Vertices(), "Checking vertices")
		})

		t.Run("should return error for no vertex", func(t *testing.T) {
			g := createGraph()

			err := g.Add()

			assert.NotNil(t, err)
		})
	})

	t.Run("Connect", func(t *testing.T) {
		t.Run("should add a new edge from root to leaf", func(t *testing.T) {
			g := createGraph()

			err := g.Connect("A", "F")

			assert.Nil(t, err)

			next, err := g.Next("A")

			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"B", "D", "F"}, next, fmt.Sprintf("Checking next of A: %v", next))

			prev, err := g.Prev("F")
			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"A", "E"}, prev, fmt.Sprintf("Checking prev of F: %v", prev))
		})

		t.Run("should return error while trying to connect existing edge", func(t *testing.T) {
			g := createGraph()

			err := g.Connect("A", "B")

			assert.NotNil(t, err)
		})

		t.Run("should return error for cycle edges", func(t *testing.T) {
			g := createGraph()

			err := g.Connect("A", "F")

			assert.Nil(t, err)

			err = g.Connect("F", "A")

			assert.NotNil(t, err)
		})

		t.Run("should return error non existing vertices", func(t *testing.T) {
			g := createGraph()

			err := g.Connect("X", "A")

			assert.NotNil(t, err)
		})
	})

	t.Run("DisconnectEdge", func(t *testing.T) {
		t.Run("should remove edge from graph", func(t *testing.T) {
			g := createGraph()

			err := g.DisconnectEdge("A", "B")

			assert.Nil(t, err)

			next, err := g.Next("A")
			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"D"}, next, "Checking next of A")

			prev, err := g.Prev("B")
			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{}, prev, "Checking prev of B")
		})

		t.Run("should remove edge of vertex with multiple prev", func(t *testing.T) {
			g := createGraph()

			err := g.DisconnectEdge("B", "E")

			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"A", "B", "C", "D", "E", "F"}, g.Vertices(), "Checking vertices")

			next, err := g.Next("B")
			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"C"}, next, "Checking next of B")
		})

		t.Run("should return error while removing non existing edge", func(t *testing.T) {
			g := createGraph()

			err := g.DisconnectEdge("A", "X")

			assert.NotNil(t, err)
		})
	})

	t.Run("Disconnect", func(t *testing.T) {
		t.Run("should disconnect a vertex", func(t *testing.T) {
			g := createGraph()

			err := g.Disconnect("B")

			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"A", "B", "C", "D", "E", "F"}, g.Vertices(), "Checking vertices")

			next, err := g.Next("B")
			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{}, next, "Checking next of A")

			prev, err := g.Prev("B")
			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{}, prev, "Checking prev of B")
		})

		t.Run("should return error for non existing vertex", func(t *testing.T) {
			g := createGraph()

			err := g.Disconnect("X")

			assert.NotNil(t, err)
		})

		t.Run("should disconnect vertex with empty prev", func(t *testing.T) {
			g := createGraph()

			err := g.Disconnect("A")
			assert.Nil(t, err)

			next, err := g.Next("A")
			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{}, next, "Checking next of A")
		})
	})

	t.Run("Remove", func(t *testing.T) {
		t.Run("should remove root vertex", func(t *testing.T) {
			g := createGraph()

			removed, err := g.Remove("A")

			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex(nil), g.Vertices(), "Checking vertices")
			assert.Equal(t, []dag.Vertex{"A", "D", "E", "F", "B", "C"}, removed, "Checking removed vertices")
		})

		t.Run("should remove middle vertex", func(t *testing.T) {
			g := createGraph()

			removed, err := g.Remove("D")

			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"A", "B", "C"}, g.Vertices(), "Checking vertices")
			assert.Equal(t, []dag.Vertex{"D", "E", "F"}, removed, "Checking removed vertices")
		})

		t.Run("should remove leaf vertex", func(t *testing.T) {
			g := createGraph()

			removed, err := g.Remove("F")

			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"A", "B", "C", "D", "E"}, g.Vertices(), "Checking vertices")
			assert.Equal(t, []dag.Vertex{"F"}, removed, "Checking removed vertices")
		})

		t.Run("should return error for non existing vertex", func(t *testing.T) {
			g := createGraph()

			removed, err := g.Remove("X")

			assert.NotNil(t, err)
			assert.Equal(t, []dag.Vertex{"A", "B", "C", "D", "E", "F"}, g.Vertices(), "Checking vertices")
			assert.Equal(t, []dag.Vertex(nil), removed, "Checking removed vertices")
		})
	})

	t.Run("TopSort", func(t *testing.T) {
		t.Run("should sort vertices in topological order", func(t *testing.T) {
			g := createGraph()

			sorted, err := g.TopSort()

			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"A", "B", "D", "C", "E", "F"}, sorted, "Checking sorted vertices")
		})

		t.Run("should sort vertices in topological order (harder case)", func(t *testing.T) {
			g := dag.New()

			var err error
			err = g.Add("2", "3", "5", "7", "8", "9", "10", "11")
			assert.Nil(t, err)

			err = g.Connect("3", "8")
			assert.Nil(t, err)
			err = g.Connect("3", "10")
			assert.Nil(t, err)

			err = g.Connect("5", "11")
			assert.Nil(t, err)

			err = g.Connect("7", "8")
			assert.Nil(t, err)
			err = g.Connect("7", "11")
			assert.Nil(t, err)

			err = g.Connect("8", "9")
			assert.Nil(t, err)

			err = g.Connect("11", "9")
			assert.Nil(t, err)
			err = g.Connect("11", "10")
			assert.Nil(t, err)

			sorted, err := g.TopSort()
			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"2", "3", "5", "7", "8", "11", "9", "10"}, sorted, "Checking sorted vertices")
		})
	})

	t.Run("DeepCopy", func(t *testing.T) {
		t.Run("should return a deep copy of a graph", func(t *testing.T) {
			g := createGraph()

			copy, err := g.DeepCopy()

			assert.Nil(t, err)

			assert.Equal(t, g.Vertices(), copy.Vertices(), "Checking vertices")
			assert.Equal(t, g.Edges(), copy.Edges(), "Checking edges")

			err = g.Append("X", []dag.Vertex{"F"})

			assert.Nil(t, err)
			assert.NotEqual(t, g.Vertices(), copy.Vertices(), "Checking vertices")
			assert.NotEqual(t, g.Edges(), copy.Edges(), "Checking edges")
		})
	})

	t.Run("SubGraph", func(t *testing.T) {
		t.Run("should return a sub graph of a graph given vertices", func(t *testing.T) {
			sub, err := createGraph().SubGraph([]dag.Vertex{"A", "B", "C"})

			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"A", "B", "C"}, sub.Vertices(), "Checking sub graph vertices")
			assert.Equal(t, dag.Edges{"A": []dag.Vertex{"B"}, "B": []dag.Vertex{"C"}, "C": []dag.Vertex{}}, sub.Edges(), "Checking sub graph edges")
		})
	})

}
