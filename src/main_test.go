package main

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/askiada/GraphDensityCut/src/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func AddEdge(gr []*model.Node, from, to int) []*model.Node {

	fromIdx := from - 1
	toIdx := to - 1
	gr[fromIdx].Neighbors = append(gr[fromIdx].Neighbors, &model.Edge{To: toIdx, Weight: 1})
	gr[toIdx].Neighbors = append(gr[toIdx].Neighbors, &model.Edge{To: fromIdx, Weight: 1})
	return gr
}

func (suite *DcutTestSuite) TestGraphOneNodeNoEdge() {
	G := make([]*model.Node, 1)
	G[0] = &model.Node{Value: 0}

	err := suite.sesh.DensityConnectedTree(G, nil)
	assert.Nil(suite.T(), err)
	minFrom, minTo, minDcut := suite.sesh.Dcut()
	assert.Equal(suite.T(), -1, minFrom)
	assert.Equal(suite.T(), -1, minTo)
	assert.Equal(suite.T(), math.Inf(1), minDcut)
}
func (suite *DcutTestSuite) TestGraphOneNodeOneEdge() {
	G := make([]*model.Node, 1)
	G[0] = &model.Node{Value: 0}
	G[0].Neighbors = append(G[0].Neighbors, &model.Edge{To: 5, Weight: 1})
	err := suite.sesh.DensityConnectedTree(G, nil)
	assert.Nil(suite.T(), err)

	minFrom, minTo, minDcut := suite.sesh.Dcut()
	assert.Equal(suite.T(), -1, minFrom)
	assert.Equal(suite.T(), -1, minTo)
	assert.Equal(suite.T(), math.Inf(1), minDcut)
}

func (suite *DcutTestSuite) TestGraphTwoNodesOneValidEdge() {
	G := GenerateGraph(2, 1)
	err := suite.sesh.DensityConnectedTree(G, nil)
	assert.Nil(suite.T(), err)
	minFrom, minTo, minDcut := suite.sesh.Dcut()
	assert.Equal(suite.T(), 1, minFrom)
	assert.Equal(suite.T(), 0, minTo)
	assert.Equal(suite.T(), 0.8306733524230347, minDcut)
}

func (suite *DcutTestSuite) TestGraphTwoNodesOneInvalidEdge() {
	G := make([]*model.Node, 2)
	G[0] = &model.Node{Value: 0}
	G[1] = &model.Node{Value: 1}
	G[0].Neighbors = append(G[0].Neighbors, &model.Edge{To: 5, Weight: 1})
	err := suite.sesh.DensityConnectedTree(G, nil)
	assert.Error(suite.T(), err)
}

func (suite *DcutTestSuite) TestGraphTwoNodesNoEdge() {
	G := make([]*model.Node, 2)
	G[0] = &model.Node{Value: 0}
	G[1] = &model.Node{Value: 1}
	err := suite.sesh.DensityConnectedTree(G, nil)
	assert.Error(suite.T(), err)
}

func (suite *DcutTestSuite) TestGraph6Nodes5Edges() {
	//1-----2-----4----6      1-----2     4----6
	//      |     |       =>        |     |
	//      3     5                 3     5
	G := make([]*model.Node, 6)
	G[0] = &model.Node{Value: 0}
	G[1] = &model.Node{Value: 1}
	G[2] = &model.Node{Value: 2}
	G[3] = &model.Node{Value: 3}
	G[4] = &model.Node{Value: 4}
	G[5] = &model.Node{Value: 5}

	AddEdge(G, 1, 2)
	AddEdge(G, 2, 3)
	AddEdge(G, 2, 4)
	AddEdge(G, 4, 5)
	AddEdge(G, 4, 6)
	err := suite.sesh.DensityConnectedTree(G, nil)
	assert.Nil(suite.T(), err)
	fmt.Println(suite.sesh.T)
	minFrom, minTo, minDcut := suite.sesh.Dcut()
	assert.Equal(suite.T(), 1, minFrom)
	assert.Equal(suite.T(), 3, minTo)
	assert.Equal(suite.T(), 0.16666666666666666, minDcut)
}

func CreateZacharyKarateClub() []*model.Node {

	graph := make([]*model.Node, 34)

	for i := 0; i < 34; i++ {
		graph[i] = &model.Node{Value: i}
	}

	graph = AddEdge(graph, 2, 1)
	graph = AddEdge(graph, 3, 1)
	graph = AddEdge(graph, 3, 2)
	graph = AddEdge(graph, 4, 1)
	graph = AddEdge(graph, 4, 2)
	graph = AddEdge(graph, 4, 3)
	graph = AddEdge(graph, 5, 1)
	graph = AddEdge(graph, 6, 1)
	graph = AddEdge(graph, 7, 1)
	graph = AddEdge(graph, 7, 5)
	graph = AddEdge(graph, 7, 6)
	graph = AddEdge(graph, 8, 1)
	graph = AddEdge(graph, 8, 2)
	graph = AddEdge(graph, 8, 3)
	graph = AddEdge(graph, 8, 4)
	graph = AddEdge(graph, 9, 1)
	graph = AddEdge(graph, 9, 3)
	graph = AddEdge(graph, 10, 3)
	graph = AddEdge(graph, 11, 1)
	graph = AddEdge(graph, 11, 5)
	graph = AddEdge(graph, 11, 6)
	graph = AddEdge(graph, 12, 1)
	graph = AddEdge(graph, 13, 1)
	graph = AddEdge(graph, 13, 4)
	graph = AddEdge(graph, 14, 1)
	graph = AddEdge(graph, 14, 2)
	graph = AddEdge(graph, 14, 3)
	graph = AddEdge(graph, 14, 4)
	graph = AddEdge(graph, 17, 6)
	graph = AddEdge(graph, 17, 7)
	graph = AddEdge(graph, 18, 1)
	graph = AddEdge(graph, 18, 2)
	graph = AddEdge(graph, 20, 1)
	graph = AddEdge(graph, 20, 2)
	graph = AddEdge(graph, 22, 1)
	graph = AddEdge(graph, 22, 2)
	graph = AddEdge(graph, 26, 24)
	graph = AddEdge(graph, 26, 25)
	graph = AddEdge(graph, 28, 3)
	graph = AddEdge(graph, 28, 24)
	graph = AddEdge(graph, 28, 25)
	graph = AddEdge(graph, 29, 3)
	graph = AddEdge(graph, 30, 24)
	graph = AddEdge(graph, 30, 27)
	graph = AddEdge(graph, 31, 2)
	graph = AddEdge(graph, 31, 9)
	graph = AddEdge(graph, 32, 1)
	graph = AddEdge(graph, 32, 25)
	graph = AddEdge(graph, 32, 26)
	graph = AddEdge(graph, 32, 29)
	graph = AddEdge(graph, 33, 3)
	graph = AddEdge(graph, 33, 9)
	graph = AddEdge(graph, 33, 15)
	graph = AddEdge(graph, 33, 16)
	graph = AddEdge(graph, 33, 19)
	graph = AddEdge(graph, 33, 21)
	graph = AddEdge(graph, 33, 23)
	graph = AddEdge(graph, 33, 24)
	graph = AddEdge(graph, 33, 30)
	graph = AddEdge(graph, 33, 31)
	graph = AddEdge(graph, 33, 32)
	graph = AddEdge(graph, 34, 9)
	graph = AddEdge(graph, 34, 10)
	graph = AddEdge(graph, 34, 14)
	graph = AddEdge(graph, 34, 15)
	graph = AddEdge(graph, 34, 16)
	graph = AddEdge(graph, 34, 19)
	graph = AddEdge(graph, 34, 20)
	graph = AddEdge(graph, 34, 21)
	graph = AddEdge(graph, 34, 23)
	graph = AddEdge(graph, 34, 24)
	graph = AddEdge(graph, 34, 27)
	graph = AddEdge(graph, 34, 28)
	graph = AddEdge(graph, 34, 29)
	graph = AddEdge(graph, 34, 30)
	graph = AddEdge(graph, 34, 31)
	graph = AddEdge(graph, 34, 32)
	graph = AddEdge(graph, 34, 33)

	return graph
}

type DcutTestSuite struct {
	suite.Suite
	G           []*model.Node
	sesh        *Session
	benchGraphs map[int]map[int][]*model.Node
}

func (suite *DcutTestSuite) SetupSuite() {
	rand.Seed(165165416)
}

func (suite *DcutTestSuite) SetupTest() {
	suite.sesh = &Session{}

}

func (suite *DcutTestSuite) TestZacharyGraph() {
	first := 7
	G := CreateZacharyKarateClub()
	err := suite.sesh.DensityConnectedTree(G, &first)
	assert.Nil(suite.T(), err)
	/*fmt.Println("T contains:", len(suite.sesh.T))
	for _, node := range suite.sesh.T {
		if node.Connect != nil {
			fmt.Println(node.Value+1, "-->", node.Connect.Value+1, node.Density, len(suite.sesh.DCTEdges[node.Value]))
		} else {
			fmt.Println(node.Value+1, "-->", "nil", node.Density, len(suite.sesh.DCTEdges[node.Value]))
		}
	}*/
	minFrom, minTo, minDcut := suite.sesh.Dcut()
	assert.Equal(suite.T(), 8, minFrom)
	assert.Equal(suite.T(), 2, minTo)
	assert.Equal(suite.T(), 0.021390374331550804, minDcut)
	//fmt.Println("Score", minFrom+1, minTo+1, minDcut)
}

/*
func (suite *DcutTestSuite) TestRandomGraph() {
	first := 4
	gra := GenerateGraph(300, 1000)
	err := suite.sesh.DensityConnectedTree(gra, &first)
	assert.Nil(suite.T(), err)
	fmt.Println("T contains:", len(T))
	for _, node := range T {
		if node.Connect != nil {
			fmt.Println(node.Value+1, "-->", node.Connect.Value+1, node.Density, len(suite.sesh.DCTEdges[node.Value]))
		} else {
			fmt.Println(node.Value+1, "-->", "nil", len(suite.sesh.DCTEdges[node.Value]))
		}

		//fmt.Println("Connect:", node.Connect)
	}
	minFrom, minTo, minDcut := suite.sesh.Dcut()
	fmt.Println("Score", minFrom+1, minTo+1, minDcut)
}
*/

func TestDcutTestSuite(t *testing.T) {
	suite.Run(t, new(DcutTestSuite))
}
