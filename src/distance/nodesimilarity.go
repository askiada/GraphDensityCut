package metric

import "github.com/askiada/Dcut/src/model"

//Definition 4
//s(u, v) = ρ(u, v) ∗ w(u, v)
func NodeSim(u *model.Node, v *model.Node, w float64) float64 {
	return w * JaccardCoeff(u, v)
}
