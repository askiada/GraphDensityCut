package metric

import (
	"github.com/askiada/GraphDensityCut/src/model"
)

var hashMapNeighStorage map[int]map[int]bool
var jaccMapResults map[int]map[int]float64

func Init(size int) {
	hashMapNeighStorage = make(map[int]map[int]bool, size)
}

//O(len(a)+len(b) * x) where x is a factor of hash function efficiency (between 1 and 2)
func countIntersect(a *model.Node, b *model.Node) float64 {
	var hashA map[int]bool
	//Γ(u) = {v ∈ V |{u, v} ∈ E} ∪ {u}
	//It means that {u,v} c Γ(u) ∩ Γ(v)
	count := float64(2)

	if _, ok := hashMapNeighStorage[a.Value]; !ok {
		hashA = make(map[int]bool)

		// len(a)
		for _, va := range a.Neighbors {
			hashA[va.To] = true
		}
		hashMapNeighStorage[a.Value] = hashA
	} else {
		hashA = hashMapNeighStorage[a.Value]
	}
	// len(b)
	for _, vb := range b.Neighbors {
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
	inter := countIntersect(a, b)
	coeff := inter / (float64(len(a.Neighbors)+len(b.Neighbors)) - inter)
	return coeff
}
