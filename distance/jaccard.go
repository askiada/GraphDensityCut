//Package metric defines all the metrics to build the density connected tree
package metric

import (
	"github.com/askiada/GraphDensityCut/model"
)

var hashMapNeighStorage map[int]map[int]bool

func Init(size int) {
	hashMapNeighStorage = make(map[int]map[int]bool, size)
}

//countIntersect Returns the number of common neighbors for two given nodes.
//Complexity: O(len(a)+len(b) * x) where x is a factor of hash function efficiency (between 1 and 2)
func countIntersect(a *model.Node, b *model.Node) float64 {
	var hashA map[int]bool
	//Γ(u) = {v ∈ V |{u, v} ∈ E} ∪ {u}
	//It means that {u,v} c Γ(u) ∩ Γ(v)
	count := float64(2)

	if _, ok := hashMapNeighStorage[a.Index]; !ok {
		hashA = make(map[int]bool)

		// len(a)
		for _, va := range a.Neighbors {
			hashA[va.To] = true
		}
		hashMapNeighStorage[a.Index] = hashA
	} else {
		hashA = hashMapNeighStorage[a.Index]
	}
	// len(b)
	for _, vb := range b.Neighbors {
		if _, ok := hashA[vb.To]; ok {
			count++
		}
	}

	return count
}

//JaccardCoeff Comput the Jaccard index based on the standard definition https://en.wikipedia.org/wiki/Jaccard_index
//
//Implements Definition 3
//  ρ(u, v) = |Γ(u) ∩ Γ(v)| / |Γ(u) ∪ Γ(v)|
//
//  |intersecrion(A,B)| / |A| + |B| - intersecrion(A,B)
func JaccardCoeff(a *model.Node, b *model.Node) float64 {
	inter := countIntersect(a, b)
	union := float64(len(a.Neighbors)+len(b.Neighbors)) + 2 - inter
	if union == 0 {
		return 1
	}
	coeff := inter / union
	return coeff
}
