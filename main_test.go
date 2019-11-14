package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"testing"

	"github.com/askiada/GraphDensityCut/model"
	"github.com/askiada/GraphDensityCut/session"
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
	G[0] = &model.Node{Value: "1", Index: 0}

	err := suite.sesh.DensityConnectedTree(G, nil)
	assert.Nil(suite.T(), err)
	minFrom, minTo, minDcut := suite.sesh.Dcut()
	assert.Equal(suite.T(), -1, minFrom)
	assert.Equal(suite.T(), -1, minTo)
	assert.Equal(suite.T(), math.Inf(1), minDcut)
}
func (suite *DcutTestSuite) TestGraphOneNodeOneEdge() {
	G := make([]*model.Node, 1)
	G[0] = &model.Node{Value: "1", Index: 0}
	G[0].Neighbors = append(G[0].Neighbors, &model.Edge{To: 5, Weight: 1})
	err := suite.sesh.DensityConnectedTree(G, nil)
	assert.Nil(suite.T(), err)

	minFrom, minTo, minDcut := suite.sesh.Dcut()
	assert.Equal(suite.T(), -1, minFrom)
	assert.Equal(suite.T(), -1, minTo)
	assert.Equal(suite.T(), math.Inf(1), minDcut)
}

func (suite *DcutTestSuite) TestGraphTwoNodesOneValidEdge() {
	G := make([]*model.Node, 2)
	G[0] = &model.Node{Value: "1", Index: 0}
	G[1] = &model.Node{Value: "2", Index: 1}
	//fmt.Println(G)
	G = AddEdge(G, 1, 2)
	/*G[0].Neighbors = append(G[0].Neighbors, &model.Edge{To: 1, Weight: 1})
	G[1].Neighbors = append(G[0].Neighbors, &model.Edge{To: 1, Weight: 1})*/
	err := suite.sesh.DensityConnectedTree(G, nil)
	assert.Nil(suite.T(), err)
	minFrom, minTo, minDcut := suite.sesh.Dcut()
	assert.Equal(suite.T(), 0, minFrom)
	assert.Equal(suite.T(), 1, minTo)
	assert.Equal(suite.T(), float64(1), minDcut)
}

func (suite *DcutTestSuite) TestGraphTwoNodesOneInvalidEdge() {
	G := make([]*model.Node, 2)
	G[0] = &model.Node{Value: "1", Index: 0}
	G[1] = &model.Node{Value: "2", Index: 1}
	G[0].Neighbors = append(G[0].Neighbors, &model.Edge{To: 5, Weight: 1})
	err := suite.sesh.DensityConnectedTree(G, nil)
	assert.Error(suite.T(), err)
}

func (suite *DcutTestSuite) TestGraphTwoNodesNoEdge() {
	G := make([]*model.Node, 2)
	G[0] = &model.Node{Value: "1", Index: 0}
	G[1] = &model.Node{Value: "2", Index: 1}
	err := suite.sesh.DensityConnectedTree(G, nil)
	assert.Error(suite.T(), err)
}

func (suite *DcutTestSuite) TestGraph6Nodes5Edges() {
	//1-----2-----4----6      1-----2     4----6      1     2     4----6
	//      |     |       =>        |     |                 |
	//      3     5                 3     5                 3     5
	G := make([]*model.Node, 6)
	G[0] = &model.Node{Value: "1", Index: 0}
	G[1] = &model.Node{Value: "2", Index: 1}
	G[2] = &model.Node{Value: "3", Index: 2}
	G[3] = &model.Node{Value: "4", Index: 3}
	G[4] = &model.Node{Value: "5", Index: 4}
	G[5] = &model.Node{Value: "6", Index: 5}
	AddEdge(G, 1, 2)
	AddEdge(G, 2, 3)
	AddEdge(G, 2, 4)
	AddEdge(G, 4, 5)
	AddEdge(G, 4, 6)
	tmp := 4
	err := suite.sesh.DensityConnectedTree(G, &tmp)
	assert.Nil(suite.T(), err)
	minFrom, minTo, minDcut := suite.sesh.Dcut()
	assert.Equal(suite.T(), "2", G[minFrom].Value)
	assert.Equal(suite.T(), "4", G[minTo].Value)
	assert.Equal(suite.T(), 0.1111111111111111, minDcut)
	p1, p2 := suite.sesh.SplitGraph()

	tmp = 2
	err = suite.sesh.DensityConnectedTree(p1, &tmp)
	assert.Nil(suite.T(), err)
	minFrom, minTo, minDcut = suite.sesh.Dcut()

	assert.Equal(suite.T(), "1", p1[minFrom].Value)
	assert.Equal(suite.T(), "2", p1[minTo].Value)
	assert.Equal(suite.T(), 0.6666666666666666, minDcut)

	p11, p12 := suite.sesh.SplitGraph()

	tmp2 := 0
	err = suite.sesh.DensityConnectedTree(p12, &tmp2)
	assert.Nil(suite.T(), err)
	minFrom, minTo, minDcut = suite.sesh.Dcut()
	assert.Equal(suite.T(), "3", p12[minFrom].Value)
	assert.Equal(suite.T(), "2", p12[minTo].Value)
	assert.Equal(suite.T(), 1.0, minDcut)

	err = suite.sesh.DensityConnectedTree(p11, &tmp2)
	assert.Nil(suite.T(), err)
	minFrom, minTo, minDcut = suite.sesh.Dcut()
	assert.Equal(suite.T(), -1, minFrom)
	assert.Equal(suite.T(), -1, minTo)
	assert.Equal(suite.T(), math.Inf(1), minDcut)

	err = suite.sesh.DensityConnectedTree(p2, &tmp)
	assert.Nil(suite.T(), err)
	minFrom, minTo, minDcut = suite.sesh.Dcut()
	assert.Equal(suite.T(), "5", p2[minFrom].Value)
	assert.Equal(suite.T(), "4", p2[minTo].Value)
	assert.Equal(suite.T(), 0.6666666666666666, minDcut)

	p21, p22 := suite.sesh.SplitGraph()

	fmt.Println(p21)
	fmt.Println(p22)

	tmp2 = 0
	err = suite.sesh.DensityConnectedTree(p22, &tmp2)
	assert.Nil(suite.T(), err)
	minFrom, minTo, minDcut = suite.sesh.Dcut()
	assert.Equal(suite.T(), "6", p22[minFrom].Value)
	assert.Equal(suite.T(), "4", p22[minTo].Value)
	assert.Equal(suite.T(), 1.0, minDcut)

	err = suite.sesh.DensityConnectedTree(p21, &tmp2)
	assert.Nil(suite.T(), err)
	minFrom, minTo, minDcut = suite.sesh.Dcut()
	assert.Equal(suite.T(), -1, minFrom)
	assert.Equal(suite.T(), -1, minTo)
	assert.Equal(suite.T(), math.Inf(1), minDcut)
}

func CreateZacharyKarateClub() []*model.Node {

	graph := make([]*model.Node, 34)

	for i := 0; i < 34; i++ {
		graph[i] = &model.Node{Value: strconv.Itoa(i + 1), Index: i}
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
	G    []*model.Node
	sesh *session.Session
}

func (suite *DcutTestSuite) SetupSuite() {
	rand.Seed(67867867)
}

func (suite *DcutTestSuite) SetupTest() {
	suite.sesh = &session.Session{}

}

func (suite *DcutTestSuite) TestZacharyGraph() {
	first := 7
	G := CreateZacharyKarateClub()
	err := suite.sesh.DensityConnectedTree(G, &first)
	assert.Nil(suite.T(), err)
	minFrom, minTo, minDcut := suite.sesh.Dcut()
	assert.Equal(suite.T(), "9", G[minFrom].Value)
	assert.Equal(suite.T(), "3", G[minTo].Value)
	assert.Equal(suite.T(), 0.01809954751131222, minDcut)

	suite.sesh.SplitGraph()
}

func TestDcutTestSuite(t *testing.T) {
	suite.Run(t, new(DcutTestSuite))
}
