package model

import "strconv"

type Edge struct {
	To     int
	Weight float64
	Check  bool
}

func (e *Edge) String() string {
	return "To: " + strconv.Itoa(e.To+1) + " Weight: " + strconv.FormatFloat(e.Weight, 'f', -1, 64)
}

type Node struct {
	Value     int
	Checked   bool
	Connect   *Node
	Density   float64
	Neighbors []*Edge
}

func (n *Node) String() string {
	neigh := "\n"

	for _, e := range n.Neighbors {
		neigh += "\t " + e.String() + "\n"
	}

	return "Node: " + strconv.Itoa(n.Value+1) + neigh
}
