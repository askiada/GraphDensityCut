package main

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/askiada/GraphDensityCut/src/graph"
	"github.com/askiada/GraphDensityCut/src/model"
	"github.com/askiada/GraphDensityCut/src/session"
)

func ResetG(g []*model.Node) {
	for i := range g {
		g[i].Checked = false
		for j := range g[i].Neighbors {
			g[i].Neighbors[j].Check = false
		}
	}
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

var Result []*model.Node

func benchmarkDcut(nodes, edges int, b *testing.B) {
	G := graph.Generate(nodes, edges)
	b.ResetTimer()
	sesh := &session.Session{}
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
