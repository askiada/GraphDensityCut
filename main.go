package main

import (
	"log"

	"github.com/askiada/GraphDensityCut/graph"
	"github.com/askiada/GraphDensityCut/session"

	"github.com/askiada/GraphDensityCut/model"
)

func run(sesh *session.Session, G []*model.Node, split bool) ([]*model.Node, []*model.Node) {
	log.Println("Build density connected tree...")
	err := sesh.DensityConnectedTree(G, nil)
	if err != nil {
		panic(err)
	}
	log.Println("Find the minimum denstity score...")
	minFrom, minTo, minDcut := sesh.Dcut()
	log.Println("Min From:", sesh.Graph[minFrom].Value)
	log.Println("Min To:", sesh.Graph[minTo].Value)
	log.Println("Min Dcut Score:", minDcut)
	if split {
		p1, p2 := sesh.SplitGraph()
		return p1, p2
	}
	return nil, nil
}

func main() {

	maxNodes := 400
	maxEdges := 45900
	G := graph.Generate(maxNodes, maxEdges)

	sesh := &session.Session{}
	p1, p2 := run(sesh, G, true)

	log.Println("Parition 1...")
	run(sesh, p1, true)
	log.Println("Parition 2...")
	run(sesh, p2, true)
}
