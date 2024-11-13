![Build Status](https://github.com/aacanakin/dag/actions/workflows/go.yml/badge.svg) [![codecov](https://codecov.io/gh/aacanakin/dag/graph/badge.svg?token=VSRDOJPW7C)](https://codecov.io/gh/aacanakin/dag) [![Go Reference](https://pkg.go.dev/badge/github.com/aacanakin/dag.svg)](https://pkg.go.dev/github.com/aacanakin/dag)

# dag

dag is an implementation of a directed acyclic graph in Go. It is a simple and easy-to-use package that allows you to create and manipulate directed acyclic graphs.

## Installation

```bash
go get github.com/aacanakin/dag
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/aacanakin/dag"
)

/*
A -> B -> C
|    |
v    v
D -> E -> F
*/
func main() {
	// Let's create a directed acyclic graph that looks like the following:
	//
	// A -> B -> C
	// |    |
	// v    v
	// D -> E -> F

	// Create an empty directed acyclic graph
	var err error
	g, err := dag.New(
		dag.WithVertices([]dag.Vertex{"A", "B", "C", "D", "E", "F"}),
		dag.WithEdges(dag.Edges{
			"A": []dag.Vertex{"B", "D"},
			"B": []dag.Vertex{"C", "E"},
			"D": []dag.Vertex{"E"},
			"E": []dag.Vertex{"F"},
		}),
	)

	if err != nil {
		panic(errors.Wrap(err, "could not create sample graph"))
	}

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
```

More examples can be found in the godoc examples.

## Roadmap
- [ ] Generic vertex types
