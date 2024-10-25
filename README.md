# dag

[![Go build](https://github.com/aacanakin/dag/actions/workflows/go.yml/badge.svg)](https://github.com/aacanakin/dag/actions/workflows/go.yml)

Directed asyclic graph implementation in go

## Install

```go
go get github.com/aacanakin/dag
```

# Usage

```go
graph := dag.New() // creates an empty directed graph

root := dag.NewNode("root")
graph.Add(root) // adds a root node

// node can also have data
root := dag.NewNode(
    "root",
    dag.WithNodeData(map[string]interface{}{ "description": "this is the root node"}),
)

mid := dag.NewNode("mid")
err := graph.Append(mid, root) // appends a node TO root node

leaf := dag.NewNode("leaf")
err := graph.Append(leaf, mid) // appends a node TO mid node

leafTree := graph.ParentTree(leaf) // returns parent tree in bfs manner
// []*dag.Node{ leaf, mid, root }

rootTree := graph.ChildTree(root) // returns child tree in bfs manner
// []*dag.Node{ root, mid, leaf }

err := graph.Connect(root, leaf) // connects root node to leaf node
// root -> mid -> leaf
// root -> leaf

nextNodes := graph.Next(root) // returns child nodes
// []*dag.Node{ mid, leaf }

prevNodes := graph.Prev(mid) // return parent nodes
// []*dag.Node{ root }
```
