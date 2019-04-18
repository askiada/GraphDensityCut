package model

import "strconv"

type Edge struct {
	//Index of the nide in the graph
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
	Index     int
	Value     string
	Checked   bool
	Connect   *Node
	Density   float64
	Neighbors []*Edge
}

func (n *Node) String() string {
	neigh := "\t Checked: " + strconv.FormatBool(n.Checked) + "\n"

	for _, e := range n.Neighbors {
		neigh += "\t " + e.String() + "\n"
	}

	return "Node Value: " + n.Value + " Index: " + strconv.Itoa(n.Index) + neigh
}
