//nolint:errcheck
package dag_test

import (
	"github.com/aacanakin/dag"
	"fmt"
)

func ExampleGraph() {
	// Create a new directed acyclic graph
	// A -> B -> C
	// |    |
	// v    v
	// D -> E -> F

	// Create an empty directed acyclic graph
	graph := dag.New()

	// Add some vertices
	graph.Add("A", "B", "C", "D", "E", "F")

	// Add some edges
	graph.Connect("A", "B")
	graph.Connect("B", "C")
	graph.Connect("A", "D")
	graph.Connect("D", "E")
	graph.Connect("B", "E")
	graph.Connect("E", "F")

	// Get the topological order
	sorted, _ := graph.TopSort()

	// Print the topological order
	for _, vertex := range sorted {
		fmt.Println(vertex)
	}

	// Output:
	// A
	// B
	// D
	// C
	// E
	// F
}

func ExampleGraph_Append() {
	graph := createGraph()

	var err error
	// Append some vertices
	err = graph.Append("G", []dag.Vertex{"F"})
	if err != nil {
		fmt.Println(err)
		return
	}

	err = graph.Append("H", []dag.Vertex{"G", "F"})
	if err != nil {
		fmt.Println(err)
		return
	}

	sorted, _ := graph.TopSort()

	// Print the topological order
	for _, vertex := range sorted {
		fmt.Println(vertex)
	}

	// Output:
	// A
	// B
	// D
	// C
	// E
	// F
	// G
	// H
}

func ExampleGraph_Deps() {
	// Create a new directed acyclic graph
	// A -> B -> C
	// |    |
	// v    v
	// D -> E -> F

	// Create an empty directed acyclic graph
	graph := dag.New()

	// Add some vertices
	graph.Add("A", "B", "C", "D", "E", "F")

	// Add some edges
	graph.Connect("A", "B")
	graph.Connect("B", "C")
	graph.Connect("A", "D")
	graph.Connect("D", "E")
	graph.Connect("B", "E")
	graph.Connect("E", "F")

	// Get the dependencies of vertex E
	deps, err := graph.Deps("E")

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, vertex := range deps {
		fmt.Println(vertex)
	}

	// Output:
	// A
	// B
	// D
}

func ExampleGraph_Disconnect() {
	// Create a new directed acyclic graph
	// A -> B -> C
	// |    |
	// v    v
	// D -> E -> F

	// Create an empty directed acyclic graph
	graph := dag.New()

	// Add some vertices
	graph.Add("A", "B", "C", "D", "E", "F")

	// Add some edges
	graph.Connect("A", "B")
	graph.Connect("B", "C")
	graph.Connect("A", "D")
	graph.Connect("D", "E")
	graph.Connect("B", "E")
	graph.Connect("E", "F")

	fmt.Println("disconnecting E")
	err := graph.Disconnect("E")
	if err != nil {
		fmt.Println(err)
		return
	}

	prev, err := graph.Prev("E")
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, prevVertex := range prev {
		fmt.Println(prevVertex)
	}

	next, err := graph.Next("E")
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, nextVertex := range next {
		fmt.Println(nextVertex)
	}

	// here, E is disconnected from all vertices, so, it should not have any dependencies or dependants
	// Output:
	// disconnecting E
}
