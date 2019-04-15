package metric

import "github.com/askiada/Dcut/src/model"

//O(len(a)+len(b) * x) where x is a factor of hash function efficiency (between 1 and 2)
func countIntersect(a []*model.Edge, b []*model.Edge) float64 {
	count := float64(0)
	hashA := make(map[int]bool)

	// len(a)
	for _, va := range a {
		hashA[va.To] = true
	}

	// len(b)
	for _, vb := range b {
		if _, ok := hashA[vb.To]; ok {
			count++
		}
	}
	return count
}

//Definition 3
//ρ(u, v) = |Γ(u) ∩ Γ(v)| / |Γ(u) ∪ Γ(v)|
// |intersecrion(A,B)| / |A| + |B| - intersecrion(A,B)

func JaccardCoeff(a *model.Node, b *model.Node) float64 {
	inter := countIntersect(a.Neighbors, b.Neighbors)
	coeff := inter / (float64(len(a.Neighbors)+len(b.Neighbors)) - inter)
	return coeff
}
