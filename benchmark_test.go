package main

import (
	"fmt"
	"testing"

	"github.com/askiada/GraphDensityCut/graph"
	"github.com/askiada/GraphDensityCut/model"
	"github.com/askiada/GraphDensityCut/session"
)

func ResetG(g []*model.Node) {
	for i := range g {
		g[i].Checked = false
		for j := range g[i].Neighbors {
			g[i].Neighbors[j].Check = false
		}
	}
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
		err := sesh.DensityConnectedTree(G, nil)
		if err != nil {
			panic(err)
		}
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
