package main

import (
	"log"

	"github.com/askiada/GraphDensityCut/graph"
	"github.com/askiada/GraphDensityCut/session"

	"github.com/askiada/GraphDensityCut/model"
)

func run(sesh *session.Session, G []*model.Node, split bool) ([]*model.Node, []*model.Node) {
	err := sesh.DensityConnectedTree(G, nil)
	if err != nil {
		panic(err)
	}
	minFrom, minTo, minDcut := sesh.Dcut()
	log.Println("Min From:", minFrom)
	log.Println("Min To:", minTo)
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

	run(sesh, p1, false)
	run(sesh, p2, true)

	//p1, p2 := sesh.SplitGraph()

	//Result = sesh.T
}
