//Package session defines the algorithm to find a minimum density cut in a graph
package session

import (
	"fmt"
	"log"
	"math"
	"math/rand"

	metric "github.com/askiada/GraphDensityCut/distance"
	"github.com/askiada/GraphDensityCut/model"
)

type Session struct {
	//Store the nodes in the graph
	Graph []*model.Node
	//Representation of the density connected tree
	DCTEdges map[int][]*model.Edge
	//Store the cardinality of the partition when we explore the density connected tree from a starting node and excluding one edge
	//map[FromNode][ExcludeNode]Cardinality of the partition
	ExploreResults map[int]map[int]int
	//Store in order the node we are adding to the density conneted tree
	T []*model.Node
	//Contains the two node indexes of the edge with the smallest D-cut score
	minFrom, minTo int
	//Smallest D-cut score
	minDcut float64
	//[To rework] Index of the next node available when we build a partition
	id int
}

//DensityConnectedTree Create the density connected tree starting from a give node. If the first node is not provided, it randomly picks one.
func (s *Session) DensityConnectedTree(Graph []*model.Node, first *int) error {
	//s.DCTCount = make(map[int]int)
	s.DCTEdges = make(map[int][]*model.Edge)
	gSize := len(Graph)
	//T = null;
	var T []*model.Node
	//Set ∀v ∈ V as unchecked (v.checked = false); --> Zero value for boolean
	//Randomly selected one node u ∈ V ;
	if first == nil {
		tmp := rand.Intn(len(Graph))
		first = &tmp
	}
	//Set u.checked = true;
	Graph[*first].Checked = true

	//u.connect = null, and u.density = null; --> Zero value for pointer connect. Density is equal to 0 by default, and it does not matter

	//T.insert(u);
	metric.Init(len(Graph))
	T = append(T, Graph[*first])
	for {
		//maxv = −1; p = null; q = null;
		maxv := float64(-1)
		var p, q *model.Node
		//while T.size < V.size do
		if len(T) >= gSize {
			break
		}
		//for j = 1 to Γ(u).size do
		for i := range T {
			u := T[i]
			if len(u.Neighbors) == 0 {
				return fmt.Errorf("Node with index %s does not have any neighbors", u)
			}
			for j := range u.Neighbors {
				//v = Γ(u).get(j);
				vEdge := u.Neighbors[j]
				if vEdge.To >= len(Graph) {
					return fmt.Errorf("Node with index %d does not exist", vEdge.To)
				}
				v := Graph[vEdge.To]
				//if v.checked == false then
				if !v.Checked {
					//If we have already computed the NodeSimilarity for an edge, we can use the score from the previous computation
					if vEdge.NodeSimilarity == nil {
						tmp := metric.NodeSim(u, v, vEdge.Weight)
						vEdge.NodeSimilarity = &tmp
					}
					//if s(u, v) > maxv then
					if *vEdge.NodeSimilarity > maxv {
						maxv = *vEdge.NodeSimilarity
						p = v
						q = u
					}
				}
			}
		}
		p.Checked = true
		p.Connect = q
		p.Density = maxv
		//After each iteration, we create a new edge in the Density Connected Tree.
		//Check is true, because we want to only check the Dcut bi-partition for one of the edge.
		s.DCTEdges[p.Index] = append(s.DCTEdges[p.Index], &model.Edge{To: q.Index, Weight: maxv, Check: true})
		s.DCTEdges[q.Index] = append(s.DCTEdges[q.Index], &model.Edge{To: p.Index, Weight: maxv})
		//T.insert(p);
		T = append(T, p)
	}
	s.Graph = Graph
	s.T = T
	return nil
}

//explore Returns the cardinality of the partition when we cut the edge between `node` and `exclude` in the density connected tree
func (s *Session) explore(node int, exclude int) int {
	val, ok := s.ExploreResults[node]

	if !ok {
		s.ExploreResults[node] = make(map[int]int)
	} else {
		if storedScore, ok2 := val[exclude]; ok2 {
			return storedScore
		}
	}

	//-1 to exclude
	count := len(s.DCTEdges[node]) - 1
	for _, edge := range s.DCTEdges[node] {
		if edge.To != exclude {
			count += s.explore(edge.To, node)
		}
	}
	s.ExploreResults[node][exclude] = count
	return count
}

//min Returns the minimum between two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

//Dcut Performs density in a graph and returns the indexes in the graph of the endpoints of the edge with the minimim Dcut score
func (s *Session) Dcut() (int, int, float64) {
	//For each nodeA we want to count the number of nodes in the partition if we break the edge between nodeA and nodeB
	//map[nodeA][nodeB]partitionSize
	s.ExploreResults = make(map[int]map[int]int)
	minDcut := math.Inf(1)
	minFrom := -1
	minTo := -1
	//For each edge in the Density Connected Tree, we evaluate the score of the two partitions defined after removing that edge
	dcut := float64(0)
	for node, edges := range s.DCTEdges {
		for _, e := range edges {
			if e.Check {
				//Count partiion should also include the node itself --> + 1
				countParition := s.explore(node, e.To) + 1
				//Dcut(C1, C2) = d(C1, C2)/min(|C1|, |C2|)
				dcut = e.Weight / float64(min(countParition, len(s.T)-countParition))

				if dcut < minDcut || (dcut == minDcut && ((node <= minFrom && e.To <= minTo) || (node <= minTo && e.To <= minFrom))) {
					minDcut = dcut
					minFrom = node
					minTo = e.To
				}
			}
		}
	}

	s.minFrom = minFrom
	s.minTo = minTo
	s.minDcut = minDcut
	return minFrom, minTo, minDcut
}

//extractParition Returns a map between the indexes of the nodes in the original graph and the indexes in the partition.
//
//The partition is generated based on the cut of the edge between `node` and `exclude` in the density connected tree.
func (s *Session) extractParition(partition map[int]int, node, exclude int) map[int]int {
	partition[node] = s.id

	for _, edge := range s.DCTEdges[node] {
		if edge.To != exclude {
			s.id++
			partition = s.extractParition(partition, edge.To, node)
		}
	}
	return partition
}

//CreatePartition Returns a partition of a graph based on the results of a Dcut.
func (s *Session) CreatePartition(from, exclude int) []*model.Node {
	paritionSize, ok := s.ExploreResults[from][exclude]
	//include from in the partition count
	paritionSize += 1
	if !ok {
		//If the count does not exsit, we can predict it
		paritionSize = len(s.Graph) - (s.ExploreResults[exclude][from] + 1)
	}

	log.Printf("Partition cardinality: %d (explore graph from vertex %s and ignore edge to %s)", paritionSize, s.Graph[from].Value, s.Graph[exclude].Value)

	partition1ID := make(map[int]int, paritionSize)
	partition1ID = s.extractParition(partition1ID, from, exclude)
	partition1 := make([]*model.Node, paritionSize)
	for node, idx := range partition1ID {
		partition1[idx] = &model.Node{}
		*partition1[idx] = *s.Graph[node]
		partition1[idx].Checked = false
		partition1[idx].Index = idx
		partition1[idx].Neighbors = nil
		newEdges := []*model.Edge{}
		for _, e := range s.Graph[node].Neighbors {
			if val, ok := partition1ID[e.To]; ok {
				tmp := &model.Edge{}
				*tmp = *e
				tmp.To = val
				tmp.NodeSimilarity = nil
				newEdges = append(newEdges, tmp)
			}
		}
		partition1[idx].Neighbors = newEdges
	}
	return partition1
}

func (s *Session) SplitGraph() ([]*model.Node, []*model.Node) {
	log.Println("Split Graph...")
	s.id = 0
	partition1 := s.CreatePartition(s.minFrom, s.minTo)
	//Reset index of partition
	s.id = 0
	partition2 := s.CreatePartition(s.minTo, s.minFrom)
	return partition1, partition2
}
