package dag_test

import (
	"fmt"
	"testing"

	"github.com/aacanakin/dag"
	"github.com/stretchr/testify/assert"
)

/*
A -- B -- C
|    |
D -- E -- F
*/
func createGraph() *dag.Graph {
	g := dag.New()
	g.Add("A")
	g.Add("B")
	g.Add("C")
	g.Add("D")
	g.Add("E")
	g.Add("F")

	g.Connect("A", "B")
	g.Connect("A", "D")
	g.Connect("B", "C")
	g.Connect("B", "E")
	g.Connect("D", "E")
	g.Connect("E", "F")
	return g
}

func TestGraph(t *testing.T) {
	t.Run("Vertices", func(t *testing.T) {
		t.Run("should return all vertices of graph", func (t *testing.T) {
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

		t.Run("should return empty slice for graph with no vertices", func (t *testing.T) {
			assert.Equal(t, 0, len(dag.New().Vertices()))
			assert.Equal(t, []dag.Vertex{}, dag.New().Vertices())
		})

		t.Run("should return true for existing vertex", func(t *testing.T) {
			g := createGraph()
			assert.Equal(t, true, g.Exists("D"))
		})
	})

	t.Run("Exists", func(t *testing.T) {
		t.Run("Exists", func (t *testing.T) {
			t.Run("should return false for non existing vertex", func (t *testing.T) {
				g := createGraph()
				assert.Equal(t, false, g.Exists("X"))
			})
		})
	})

	t.Run("Prev", func(t *testing.T) {
		t.Run("should return error for non existing vertex", func (t *testing.T) {
			g := createGraph()
			prev, err := g.Prev("X")
			assert.Equal(t, []dag.Vertex{}, prev)
			assert.NotNil(t, err)
		})

		t.Run("should return empty prev for root vertex", func (t *testing.T) {
			g := createGraph()
			prev, err := g.Prev("A")
			assert.Equal(t, []dag.Vertex{}, prev)
			assert.Nil(t, err)
		})

		t.Run("should return single previous vertex", func (t *testing.T) {
			g := createGraph()
			prev, err := g.Prev("B")
			assert.Equal(t, []dag.Vertex{"A"}, prev)
			assert.Nil(t, err)
		})

		t.Run("should return multiple previous vertices", func (t *testing.T) {
			g := createGraph()
			prev, err := g.Prev("E")
			assert.Equal(t, []dag.Vertex{"B", "D"}, prev)
			assert.Nil(t, err)
		})
	})

	t.Run("Next", func(t *testing.T) {
		t.Run("should return error for non existing vertex", func (t *testing.T) {
			g := createGraph()
			next, err := g.Next("X")
			assert.Equal(t, []dag.Vertex{}, next)
			assert.NotNil(t, err)
		})

		t.Run("should return empty next for leaf vertex", func (t *testing.T) {
			g := createGraph()
			next, err := g.Next("C")
			assert.Equal(t, []dag.Vertex{}, next)
			assert.Nil(t, err)
		})

		t.Run("should return single next vertex", func (t *testing.T) {
			g := createGraph()
			next, err := g.Next("D")
			assert.Equal(t, []dag.Vertex{"E"}, next)
			assert.Nil(t, err)
		})
	})

	t.Run("ReverseEdges", func(t *testing.T) {
		t.Run("should return no edges with empty edges", func (t *testing.T) {
			g := dag.New()
			g.Add(dag.Vertex("A"))

			revEdges, err := g.ReverseEdges()
			expected := dag.Edges{"A": {}}
			assert.Nil(t, err)
			assert.Equal(t, expected, revEdges)
		})

		t.Run("should reverse single edge correctly", func (t *testing.T) {
			g := dag.New()
			g.Add("A")
			g.Add("B")
			g.Connect("A", "B")

			revEdges, err := g.ReverseEdges()
			expected := dag.Edges{"A": {}, "B": {"A"}}
			assert.Nil(t, err)
			assert.Equal(t, expected, revEdges)
		})

		t.Run("should reverse graph with multiple next for a single vertex correctly", func (t *testing.T) {
			g := dag.New()
			g.Add(dag.Vertex("A"))
			g.Add(dag.Vertex("B"))
			g.Add(dag.Vertex("C"))
			g.Connect("A", "B")
			g.Connect("A", "C")

			revEdges, err := g.ReverseEdges()
			expected := dag.Edges{"A": {}, "B": {"A"}, "C": {"A"}}
			assert.Nil(t, err)
			assert.Equal(t, expected, revEdges)
		})

		t.Run("should reverse sample graph", func (t *testing.T) {
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
		t.Run("should reverse sample graph", func (t *testing.T) {
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
		t.Run("should return dfs for root vertex", func (t *testing.T) {
			g := createGraph()

			dfs, err := g.DFS("A")
			expected := []dag.Vertex{"A", "D", "E", "F", "B", "C"}

			assert.Nil(t, err)
			assert.Equal(t, expected, dfs)
		})

		t.Run("should return dfs for leaf vertex", func (t *testing.T) {
			g := createGraph()

			dfs, err := g.DFS("C")

			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"C"}, dfs)
		})

		t.Run("should return dfs for middle vertex", func (t *testing.T) {
			g := createGraph()

			dfs, err := g.DFS("B")

			assert.Nil(t, err)

			expected := []dag.Vertex{"B", "E", "F", "C"}
			assert.Equal(t, expected, dfs)
		})

		t.Run("should return error for non existing vertex", func (t *testing.T) {
			g := createGraph()

			dfs, err := g.DFS("X")

			assert.NotNil(t, err)
			assert.Equal(t, []dag.Vertex{}, dfs)
		})
	})

	t.Run("BFS", func(t *testing.T) {
		t.Run("should return bfs for root vertex", func (t *testing.T) {
			g := createGraph()

			bfs, err := g.BFS("A")

			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"A", "B", "D", "C", "E", "F"}, bfs)
		})

		t.Run("should return bfs for leaf vertex", func (t *testing.T) {
			g := createGraph()

			bfs, err := g.BFS("F")

			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"F"}, bfs)
		})

		t.Run("should return bfs for middle vertex", func (t *testing.T) {
			g := createGraph()

			bfs, err := g.BFS("B")

			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"B", "C", "E", "F"}, bfs)
		})

		t.Run("should return error for non existing vertex", func (t *testing.T) {
			g := createGraph()

			bfs, err := g.BFS("X")

			assert.NotNil(t, err)
			assert.Equal(t, []dag.Vertex{}, bfs)
		})
	})

	type depsResult  struct {
		deps []dag.Vertex
		err  error
	}

	t.Run("Deps", func(t *testing.T) {

		t.Run("should return deps for every vertex in sample graph", func (t *testing.T) {
			g := createGraph()

			results := map[string]depsResult{}
			expected := map[string]depsResult{
				"A": {[]dag.Vertex(nil), nil},
				"B": {[]dag.Vertex{"A"}, nil},
				"C": {[]dag.Vertex{"A", "B"}, nil},
				"D": {[]dag.Vertex{"A"}, nil},
				"E": {[]dag.Vertex{"A", "B", "D"}, nil},
				"F": {[]dag.Vertex{"A", "B", "D", "E"}, nil},
			}

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
		t.Run("should return reverse deps for every vertex in sample graph", func (t *testing.T) {
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
		t.Run("should return leaf vertices for sample graph", func (t *testing.T) {
			g := createGraph()

			leaves, err := g.Leaves()

			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"C", "F"}, leaves)
		})

		t.Run("should return leaf vertices for reversed sample graph", func (t *testing.T) {
			g, err := createGraph().Reverse()

			assert.Nil(t, err)

			leaves, err := g.Leaves()

			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"A"}, leaves)
		})
	})

	t.Run("Roots", func(t *testing.T) {
		t.Run("should return root vertices for sample graph", func (t *testing.T) {
			g := createGraph()

			roots, err := g.Roots()

			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"A"}, roots)
		})

		t.Run("should return root vertices for reversed sample graph", func (t *testing.T) {
			g, err := createGraph().Reverse()

			assert.Nil(t, err)

			roots, err := g.Roots()

			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"C", "F"}, roots)
		})

		t.Run("should return multi disconnected root vertices", func (t *testing.T) {
			g := dag.New()
			g.Add("A")
			g.Add("B")

			roots, err := g.Roots()
			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"A", "B"}, roots)
		})

		t.Run("should return multi interconnected root vertices", func (t *testing.T) {
			g := dag.New()
			g.Add("A")
			g.Add("B")
			g.Add("C")
			g.Add("D")

			g.Connect("A", "B")
			g.Connect("C", "D")

			roots, err := g.Roots()
			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"A", "C"}, roots)
		})

		t.Run("should return multi connected root vertices", func (t *testing.T) {
			g := dag.New()
			g.Add("R1")
			g.Add("R2")
			g.Add("A")
			g.Add("B")

			g.Connect("R1", "A")
			g.Connect("R2", "A")
			g.Connect("A", "B")

			roots, err := g.Roots()
			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"R1", "R2"}, roots)
		})
	})

	t.Run("Append", func(t *testing.T) {
		t.Run("should append new vertex to sample graph", func (t *testing.T) {
			g := createGraph()

			g.Append("X", "A", "B")
			assert.Equal(t, []dag.Vertex{"A", "B", "C", "D", "E", "F", "X"}, g.Vertices(), "Checking vertices")

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

		t.Run("should append new vertex without any prev vertice", func (t *testing.T) {
			g := createGraph()

			g.Append("X")
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

		t.Run("should append new vertex with every vertex as prev", func (t *testing.T) {
			g := createGraph()

			g.Append("X", "A", "B", "C", "D", "E", "F")

			assert.Equal(t, []dag.Vertex{"A", "B", "C", "D", "E", "F", "X"}, g.Vertices(), "Checking vertices")

			prev, err := g.Prev("X")
			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"A", "B", "C", "D", "E", "F"}, prev, fmt.Sprintf("Checking prev of X: %v", prev))

			next, err := g.Next("X")
			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{}, next, fmt.Sprintf("Checking next of X: %v", next))
		})
	})

	t.Run("Add", func(t *testing.T) {
		t.Run("should add new vertex to sample graph", func (t *testing.T) {
			g := createGraph()

			err := g.Add("X")

			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"A", "B", "C", "D", "E", "F", "X"}, g.Vertices(), "Checking vertices")
		})

		t.Run("should return error while adding existing vertex to sample graph", func (t *testing.T) {
			g := createGraph()

			err := g.Add("A")

			assert.NotNil(t, err)
			assert.Equal(t, []dag.Vertex{"A", "B", "C", "D", "E", "F"}, g.Vertices(), "Checking vertices")
		})
	})

	t.Run("Connect", func(t *testing.T) {
		t.Run("should add a new edge from root to leaf", func (t *testing.T) {
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

		t.Run("should return error while trying to connect existing edge", func (t *testing.T) {
			g := createGraph()

			err := g.Connect("A", "B")

			assert.NotNil(t, err)
		})

		t.Run("should return error for cycle edges", func (t *testing.T) {
			g := createGraph()

			err := g.Connect("A", "F")

			assert.Nil(t, err)

			err = g.Connect("F", "A")

			assert.NotNil(t, err)
		})
	})

	t.Run("DisconnectEdge", func(t *testing.T) {
		t.Run("should remove edge from graph", func (t *testing.T) {
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

		t.Run("should remove edge of vertex with multiple prev", func (t *testing.T) {
			g := createGraph()

			err := g.DisconnectEdge("B", "E")

			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"A", "B", "C", "D", "E", "F"}, g.Vertices(), "Checking vertices")

			next, err := g.Next("B")
			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"C"}, next, "Checking next of B")
		})

		t.Run("should return error while removing non existing edge", func (t *testing.T) {
			g := createGraph()

			err := g.DisconnectEdge("A", "X")

			assert.NotNil(t, err)
		})
	})

	t.Run("Disconnect", func(t *testing.T) {
		t.Run("should disconnect a vertex", func (t *testing.T) {
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

		t.Run("should return error for non existing vertex", func (t *testing.T) {
			g := createGraph()

			err := g.Disconnect("X")

			assert.NotNil(t, err)
		})

		t.Run("should disconnect vertex with empty prev", func (t *testing.T) {
			g := createGraph()

			err := g.Disconnect("A")

			next, err := g.Next("A")
			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{}, next, "Checking next of A")
		})
	})

	t.Run("Remove", func(t *testing.T) {
		t.Run("should remove root vertex", func (t *testing.T) {
			g := createGraph()

			removed, err := g.Remove("A")

			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex(nil), g.Vertices(), "Checking vertices")
			assert.Equal(t, []dag.Vertex{"A", "D", "E", "F", "B", "C"}, removed, "Checking removed vertices")
		})

		t.Run("should remove middle vertex", func (t *testing.T) {
			g := createGraph()

			removed, err := g.Remove("D")

			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"A", "B", "C"}, g.Vertices(), "Checking vertices")
			assert.Equal(t, []dag.Vertex{"D", "E", "F"}, removed, "Checking removed vertices")
		})

		t.Run("should remove leaf vertex", func (t *testing.T) {
			g := createGraph()

			removed, err := g.Remove("F")

			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"A", "B", "C", "D", "E"}, g.Vertices(), "Checking vertices")
			assert.Equal(t, []dag.Vertex{"F"}, removed, "Checking removed vertices")
		})

		t.Run("should return error for non existing vertex", func (t *testing.T) {
			g := createGraph()

			removed, err := g.Remove("X")

			assert.NotNil(t, err)
			assert.Equal(t, []dag.Vertex{"A", "B", "C", "D", "E", "F"}, g.Vertices(), "Checking vertices")
			assert.Equal(t, []dag.Vertex(nil), removed, "Checking removed vertices")
		})
	})

	t.Run("TopSort", func(t *testing.T) {
		t.Run("should sort vertices in topological order", func (t *testing.T) {
			g := createGraph()

			sorted, err := g.TopSort()

			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"A", "B", "D", "C", "E", "F"}, sorted, "Checking sorted vertices")
		})

		t.Run("should sort vertices in topological order (harder case)", func (t *testing.T) {
			g := dag.New()
			g.Add("2")
			g.Add("3")
			g.Add("5")
			g.Add("7")
			g.Add("8")
			g.Add("9")
			g.Add("10")
			g.Add("11")

			g.Connect("3", "8")
			g.Connect("3", "10")

			g.Connect("5", "11")

			g.Connect("7", "8")
			g.Connect("7", "11")

			g.Connect("8", "9")

			g.Connect("11", "9")
			g.Connect("11", "10")

			sorted, err := g.TopSort()
			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"2", "3", "5", "7", "8", "11", "9", "10"}, sorted, "Checking sorted vertices")
		})
	})

	t.Run("DeepCopy", func(t *testing.T) {
		t.Run("should return a deep copy of a graph", func (t *testing.T) {
			g := createGraph()

			copy := g.DeepCopy()

			assert.Equal(t, g.Vertices(), copy.Vertices(), "Checking vertices")
			assert.Equal(t, g.Edges(), copy.Edges(), "Checking edges")

			g.Append("X", "F")

			assert.NotEqual(t, g.Vertices(), copy.Vertices(), "Checking vertices")
			assert.NotEqual(t, g.Edges(), copy.Edges(), "Checking edges")
		})
	})

	t.Run("SubGraph", func(t *testing.T) {
		t.Run("should return a sub graph of a graph given vertices", func (t *testing.T) {
			sub, err := createGraph().SubGraph([]dag.Vertex{"A", "B", "C"})

			assert.Nil(t, err)
			assert.Equal(t, []dag.Vertex{"A", "B", "C"}, sub.Vertices(), "Checking sub graph vertices")
			assert.Equal(t, dag.Edges{"A": []dag.Vertex{"B"}, "B": []dag.Vertex{"C"}, "C": []dag.Vertex{}}, sub.Edges(), "Checking sub graph edges")
		})
	})

}
