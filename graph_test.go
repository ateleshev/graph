// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"
	"math/big"

	"github.com/soniakeys/graph"
)

func ExampleAdjacencyList_BoundsOk() {
	var g graph.AdjacencyList
	ok, _, _ := g.BoundsOk() // zero value adjacency list is valid
	fmt.Println(ok)
	g = graph.AdjacencyList{
		0: {9},
	}
	fmt.Println(g.BoundsOk()) // arc 0 to 9 invalid with only one node
	// Output:
	// true
	// false 0 9
}

func ExampleAdjacencyList_HasParallelMap() {
	g := graph.AdjacencyList{
		1: {0},
	}
	fmt.Println(g.HasParallelMap()) // result -1 -1 means no parallel arc
	g[1] = append(g[1], 0)          // add parallel arc
	// result 1 0 means parallel arc from node 1 to node 0
	fmt.Println(g.HasParallelMap())
	// Output:
	// -1 -1
	// 1 0
}

func ExampleAdjacencyList_HasParallelSort() {
	g := graph.AdjacencyList{
		1: {0},
	}
	fmt.Println(g.HasParallelSort()) // result -1 -1 means no parallel arc
	g[1] = append(g[1], 0)           // add parallel arc
	// result 1 0 means parallel arc from node 1 to node 0
	fmt.Println(g.HasParallelSort())
	// Output:
	// -1 -1
	// 1 0
}

func ExampleAdjacencyList_IsSimple() {
	// arcs directed down
	//   2
	//  / \
	// 0   1
	g := graph.AdjacencyList{
		2: {0, 1},
	}
	fmt.Println(g.IsSimple())
	// Output:
	// true -1
}

func ExampleAdjacencyList_IsSimple_loop() {
	// arcs directed down
	//   2
	//  / \
	// 0   1---\
	//      \--/
	g := graph.AdjacencyList{
		2: {0, 1},
		1: {1}, // loop
	}
	fmt.Println(g.IsSimple())
	// Output:
	// false 1
}

func ExampleAdjacencyList_IsSimple_parallelArc() {
	// arcs directed down
	//   2
	//  /|\
	//  |/ \
	//  0   1
	g := graph.AdjacencyList{
		2: {0, 1, 0},
	}
	fmt.Println(g.IsSimple())
	// Output:
	// false 2
}

func ExampleOneBits() {
	g := make(graph.AdjacencyList, 5)
	var b big.Int
	fmt.Printf("%b\n", graph.OneBits(&b, len(g)))
	// Output:
	// 11111
}
