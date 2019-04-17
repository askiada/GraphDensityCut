package main

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/askiada/GraphDensityCut/src/model"
)

func ResetG(g []*model.Node) {
	for i := range g {
		g[i].Checked = false
		for j := range g[i].Neighbors {
			g[i].Neighbors[j].Check = false
		}
	}
}

type Pair struct {
	a int
	b int
}

func GenerateEdge(maxNodes int, from int) *model.Edge {
	to := 0
	for true {
		to = rand.Intn(maxNodes)
		if from != to {
			break
		}
	}
	e := model.Edge{To: to, Weight: rand.Float64()}
	return &e
}

func GenerateGraph(maxNodes int, maxEdges int) []*model.Node {
	maxPairs := (maxNodes * (maxNodes - 1)) / 2

	nonConnectedEdges := maxEdges - maxNodes
	edgeMap := make(map[int]map[int]float64)
	s := rand.Perm(maxNodes)
	gra := make([]*model.Node, maxNodes)
	for i, node := range s {
		gra[node] = &model.Node{Value: node}
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

var Result []*model.Node

func benchmarkDcut(nodes, edges int, b *testing.B) {
	G := GenerateGraph(nodes, edges)
	b.ResetTimer()
	sesh := &Session{}
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		ResetG(G)
		b.StartTimer()
		sesh.DensityConnectedTree(G, nil)
		sesh.Dcut()
	}
	Result = sesh.T
}

func CreateEdgesCountSlice(nodesCount int) []int {
	s := []int{}
	maxEdges := (nodesCount * (nodesCount - 1)) / 2
	for i := maxEdges; i >= nodesCount; i = i / 2 {
		s = append(s, i)
	}
	return s
}

func BenchmarkDcut(b *testing.B) {
	nodesList := []int{
		10,
		25,
		100,
		200,
		400,
	}

	for _, nodesCount := range nodesList {
		edgesList := CreateEdgesCountSlice(nodesCount)
		for i := len(edgesList) - 1; i >= 0; i-- {
			b.Run(fmt.Sprintf("Nodes-%d:Edges%d", nodesCount, edgesList[i]), func(b *testing.B) {
				benchmarkDcut(nodesCount, edgesList[i], b)
			})
		}
	}
}
