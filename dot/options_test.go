package dot_test

import (
	"fmt"
	"os"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/dot"
)

func ExampleDirected() {
	// arcs directed down:
	// 0  2
	// | /|
	// |/ |
	// 3  4
	g := graph.AdjacencyList{
		0: {3},
		2: {3, 4},
		4: {},
	}
	// default for AdjacencyList is directed
	dot.WriteAdjacencyList(g, os.Stdout, dot.Indent(""))
	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout)

	// Directed(false) generates error witout reciprocal arcs
	err := dot.WriteAdjacencyList(g, os.Stdout, dot.Directed(false))
	fmt.Fprintln(os.Stdout, "Error:", err)
	fmt.Fprintln(os.Stdout)

	// undirected
	u := graph.Directed{g}.Undirected()
	dot.WriteAdjacencyList(u.AdjacencyList, os.Stdout,
		dot.Directed(false), dot.Indent(""))

	// Output:
	// digraph {
	// 0 -> 3
	// 2 -> {3 4}
	// }
	//
	// Error: directed graph
	//
	// graph {
	// 0 -- 3
	// 2 -- {3 4}
	// }
}

func ExampleEdgeLabel() {
	// arcs directed down:
	//      0       4
	// (.33)|      /|
	//      | (1.7) |
	//      |/      |(2e117)
	//      2       3
	weights := map[int]float64{
		30: .33,
		20: 1.7,
		10: 2e117,
	}
	lf := func(l graph.LI) string {
		return fmt.Sprintf(`"%g"`, weights[int(l)])
	}
	g := graph.LabeledAdjacencyList{
		0: {{2, 30}},
		4: {{2, 20}, {3, 10}},
	}
	dot.WriteLabeledAdjacencyList(g, os.Stdout,
		dot.EdgeLabel(lf), dot.Indent(""))
	// Output:
	// digraph {
	// 0 -> 2 [label = "0.33"]
	// 4 -> 2 [label = "1.7"]
	// 4 -> 3 [label = "2e+117"]
	// }
}

func ExampleGraphAttr() {
	// arcs directed right:
	// 0---2
	//  \ / \
	//   1---3
	g := graph.AdjacencyList{
		0: {1, 2},
		1: {2, 3},
		2: {3},
		3: {},
	}
	dot.WriteAdjacencyList(g, os.Stdout,
		dot.GraphAttr("rankdir", "LR"), dot.Indent(""))
	// Output:
	// digraph {
	// rankdir = LR
	// 0 -> {1 2}
	// 1 -> {2 3}
	// 2 -> 3
	// }
}

func ExampleIndent() {
	// arcs directed down:
	// 0  4
	// | /|
	// |/ |
	// 2  3
	g := graph.AdjacencyList{
		0: {2},
		4: {2, 3},
	}
	// All other examples have Indent("") to avoid a quirk of go test
	// that it can't handle leading space in the output.  In this example a
	// nonbreaking space works around the quirk to show indented output that
	// looks like the default two space indent.
	// (But then if you render it with graphviz, graphviz picks up the nbsp
	// as a node statement...)
	dot.WriteAdjacencyList(g, os.Stdout, dot.Indent("\u00a0 "))
	// Output:
	// digraph {
	//   0 -> 2
	//   4 -> {2 3}
	// }
}

func ExampleIsolated() {
	// 0  1-->2
	g := graph.AdjacencyList{
		1: {2},
		2: {},
	}
	dot.WriteAdjacencyList(g, os.Stdout, dot.Isolated(true), dot.Indent(""))
	// Output:
	// digraph {
	// 0
	// 1 -> 2
	// }
}

func ExampleNodeLabel() {
	// arcs directed down:
	// A  D
	// | /|
	// |/ |
	// B  C
	labels := []string{
		0: "A",
		4: "D",
		2: "B",
		3: "C",
	}
	lf := func(n graph.NI) string { return labels[n] }
	g := graph.AdjacencyList{
		0: {2},
		4: {2, 3},
	}
	dot.WriteAdjacencyList(g, os.Stdout, dot.Indent(""), dot.NodeLabel(lf))
	// Output:
	// digraph {
	// A -> B
	// D -> {B C}
	// }
}

func ExampleNodeLabel_construction() {
	// arcs directed down:
	// A  D
	// | /|
	// |/ |
	// B  C
	var g graph.AdjacencyList

	// example graph construction mechanism
	labels := []string{}
	nodes := map[string]graph.NI{}
	node := func(l string) graph.NI {
		if n, ok := nodes[l]; ok {
			return n
		}
		n := graph.NI(len(labels))
		labels = append(labels, l)
		g = append(g, nil)
		nodes[l] = n
		return n
	}
	addArc := func(fr, to string) {
		f := node(fr)
		g[f] = append(g[f], node(to))
	}

	// construct graph
	addArc("A", "B")
	addArc("D", "B")
	addArc("D", "C")

	// generate dot
	lf := func(n graph.NI) string { return labels[n] }
	dot.WriteAdjacencyList(g, os.Stdout, dot.Indent(""), dot.NodeLabel(lf))

	// Output:
	// digraph {
	// A -> B
	// D -> {B C}
	// }
}

func ExampleUndirectArcs() {
	//              (label 0, wt 1.6)
	//          0----------------------2
	// (label 1 |                     /
	//  wt .33) |  ------------------/
	//          | / (label 2, wt 1.7)
	//          |/
	//          1
	weights := []float64{
		0: 1.6,
		1: .33,
		2: 1.7,
	}
	g := graph.WeightedEdgeList{
		WeightFunc: func(l graph.LI) float64 { return weights[int(l)] },
		Order:      3,
		Edges: []graph.LabeledEdge{
			{graph.Edge{0, 1}, 1},
			{graph.Edge{0, 2}, 0},
			{graph.Edge{1, 2}, 2},
		},
	}
	dot.WriteWeightedEdgeList(g, os.Stdout,
		dot.UndirectArcs(true), dot.Indent(""))
	// Output:
	// graph {
	// 0 -- 1 [label = "0.33"]
	// 0 -- 2 [label = "1.6"]
	// 1 -- 2 [label = "1.7"]
	// }
}