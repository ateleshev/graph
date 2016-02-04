// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// graph.go
//
// Definitions for unlabeled graphs, and methods not specific to directed
// or undirected graphs.  Method docs need not mention that they work on both.

package graph

import (
	"math/big"
	"sort"
)

//go:generate cp cg_label.go cg_adj.go
//go:generate gofmt -r "LabeledAdjacencyList -> AdjacencyList" -w cg_adj.go
//go:generate gofmt -r "n.To -> n" -w cg_adj.go
//go:generate gofmt -r "Half -> NI" -w cg_adj.go

var one = big.NewInt(1)

// OneBits sets a big.Int to a number that is all 1s in binary.
//
// It's a utility function useful for initializing a bitmap of a graph
// to all ones; that is, with a bit set to 1 for each node of the graph.
//
// OneBits modifies b, then returns b for convenience.
func OneBits(b *big.Int, n int) *big.Int {
	return b.Sub(b.Lsh(one, uint(n)), one)
}

// NI is a "node int"
//
// It is a node number.  It is used extensively as a slice index.
// Node numbers also account for a significant fraction of the memory
// required to represent a graph.
type NI int32

type NodeList []NI

func (l NodeList) Len() int           { return len(l) }
func (l NodeList) Less(i, j int) bool { return l[i] < l[j] }
func (l NodeList) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }

// An AdjacencyList represents a graph as a list of neighbors for each node.
// The "node ID" of a node is simply it's slice index in the AdjacencyList.
//
// Adjacency lists are inherently directed. To represent an undirected graph,
// create reciprocal neighbors.
type AdjacencyList [][]NI

// HasParallelSort identifies if a graph contains parallel arcs, multiple arcs
// that lead from a node to the same node.
//
// If the graph has parallel arcs, the results fr and to represent an example
// where there are parallel arcs from node fr to node to.
//
// If there are no parallel arcs, the method returns -1 -1.
//
// Multiple loops on a node count as parallel arcs.
//
// "Sort" in the method name indicates that sorting is used to detect parallel
// arcs.  Compared to method HasParallelMap, this may give better performance
// for small or sparse graphs but will have asymtotically worse performance for
// large dense graphs.
func (g AdjacencyList) HasParallelSort() (fr, to NI) {
	var t NodeList
	for n, to := range g {
		if len(to) == 0 {
			continue
		}
		t = append(t[:0], to...)
		sort.Sort(t)
		t0 := t[0]
		for _, to := range t[1:] {
			if to == t0 {
				return NI(n), t0
			}
			t0 = to
		}
	}
	return -1, -1
}
