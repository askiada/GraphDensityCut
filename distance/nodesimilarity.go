package metric

import "github.com/askiada/GraphDensityCut/model"

//NodeSim Computes the node similarity based on the definition 4.
//
//  s(u, v) = ρ(u, v) ∗ w(u, v)
func NodeSim(u *model.Node, v *model.Node, w float64) float64 {
	return w * JaccardCoeff(u, v)
}
