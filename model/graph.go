//Package model is a warehouse for all the structures to represent a graph, especially to solve the D-cut problem
package model

import "strconv"

type Edge struct {
	//Index of the node in the graph
	To int
	//Weight of the edge
	Weight float64
	//If set to true, the Dcut will explore that edge (useful whith undirected graphs)
	Check bool
	//Store the value of the NodeSimilarity between the start and the end of the edge
	NodeSimilarity *float64
}

func (e *Edge) String() string {
	return "To Index: " + strconv.Itoa(e.To) + " Weight: " + strconv.FormatFloat(e.Weight, 'f', -1, 64) + " Check:" + strconv.FormatBool(e.Check)

}

type Node struct {
	//Index of the node in the graph
	Index int
	//Label linked to the node
	Value string
	//To ensure we only access the node one time when we build the density connected tree
	Checked bool
	//[Deprecated] Same role as checked
	Connect *Node
	//[Deprecated] Same Role as Edge.NodeSimilarity
	Density float64
	//List of neighbors
	Neighbors []*Edge
}

func (n *Node) String() string {
	neigh := "\t Checked: " + strconv.FormatBool(n.Checked) + "\n"

	for _, e := range n.Neighbors {
		neigh += "\t " + e.String() + "\n"
	}

	return "Node Value: " + n.Value + " Index: " + strconv.Itoa(n.Index) + neigh
}
