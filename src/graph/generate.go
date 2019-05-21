package graph

import (
	"math/rand"
	"strconv"

	"github.com/askiada/GraphDensityCut/src/model"
)

type Pair struct {
	a int
	b int
}

func Generate(maxNodes int, maxEdges int) []*model.Node {
	maxPairs := (maxNodes * (maxNodes - 1)) / 2

	nonConnectedEdges := maxEdges - maxNodes
	edgeMap := make(map[int]map[int]float64)
	s := rand.Perm(maxNodes)
	gra := make([]*model.Node, maxNodes)
	for i, node := range s {
		gra[node] = &model.Node{Value: strconv.Itoa(node + 1), Index: node}
		if i > 0 {
			if _, ok := edgeMap[s[i-1]]; !ok {
				edgeMap[s[i-1]] = make(map[int]float64)
			}
			if _, ok := edgeMap[node]; !ok {
				edgeMap[node] = make(map[int]float64)
			}
			w := rand.Float64()
			gra[s[i-1]].Neighbors = append(gra[s[i-1]].Neighbors, &model.Edge{To: node, Weight: w})
			gra[node].Neighbors = append(gra[node].Neighbors, &model.Edge{To: s[i-1], Weight: w})
			edgeMap[s[i-1]][node] = w
			edgeMap[node][s[i-1]] = w
		}
	}

	remainingPairs := maxPairs - (maxNodes - 1)

	pairs := make([]Pair, remainingPairs, remainingPairs)

	k := 0
	for i := 0; i < maxNodes-1; i++ {
		for j := i + 1; j < maxNodes; j++ {
			if _, ok := edgeMap[i][j]; !ok {
				//fmt.Println("(i,j)", i, j)
				pairs[k] = Pair{a: i, b: j}
				k++
			}
		}
	}

	rand.Shuffle(len(pairs), func(i, j int) { pairs[i], pairs[j] = pairs[j], pairs[i] })

	for i := 0; i < nonConnectedEdges; i++ {

		p := pairs[i]
		w := rand.Float64()
		gra[p.a].Neighbors = append(gra[p.a].Neighbors, &model.Edge{To: p.b, Weight: w})
		gra[p.b].Neighbors = append(gra[p.b].Neighbors, &model.Edge{To: p.a, Weight: w})

	}

	return gra
}
