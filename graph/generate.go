//Pacakge graph defines function to generate/update a graph
package graph

import (
	"log"
	"math/rand"
	"strconv"

	"github.com/askiada/GraphDensityCut/model"
)

type pair struct {
	a int
	b int
}

//Generate Build a random graphs with maxNodes Vertices and maxEdges Edges.
// The graph contains at least maxNodes edges to prevent isolated nodes
// Speed is not a constraint, it has a very bad time and space complexity.
//
// What are you doing ?
//
// - Counts the number of distinct pairs of nodes
//
// - Creates at least one edge per node to ensure there is no isolated nodes
//
// - Build a list of pairs with all the remaining candidates edges
//
// - Add edges to the graph until we reach maxEdges
func Generate(maxNodes int, maxEdges int) []*model.Node {
	log.Printf("Generate a random graph with %d vertices and a maximum of %d edges", maxNodes, maxEdges)
	//Number of combinations without replacement C(maxNode,2)
	maxPairs := (maxNodes * (maxNodes - 1)) / 2
	//Store the remaining edged we must create to have maxEdges and assuming we create at least one edge per node
	nonConnectedEdges := maxEdges - maxNodes
	//Store existing edges
	edgeMap := make(map[int]map[int]float64)
	s := rand.Perm(maxNodes)
	gra := make([]*model.Node, maxNodes)
	//Build at least one edge for each node
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

	//Build all the remaining edges
	remainingPairs := maxPairs - (maxNodes - 1)
	pairs := make([]pair, remainingPairs)
	k := 0
	for i := 0; i < maxNodes-1; i++ {
		for j := i + 1; j < maxNodes; j++ {
			if _, ok := edgeMap[i][j]; !ok {
				pairs[k] = pair{a: i, b: j}
				k++
			}
		}
	}

	rand.Shuffle(len(pairs), func(i, j int) { pairs[i], pairs[j] = pairs[j], pairs[i] })
	//Add edges based on the remaining pairs
	for i := 0; i < nonConnectedEdges; i++ {
		p := pairs[i]
		w := rand.Float64()
		gra[p.a].Neighbors = append(gra[p.a].Neighbors, &model.Edge{To: p.b, Weight: w})
		gra[p.b].Neighbors = append(gra[p.b].Neighbors, &model.Edge{To: p.a, Weight: w})

	}

	return gra
}
