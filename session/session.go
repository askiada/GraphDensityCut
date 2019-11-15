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

type Chain struct {
	head *NodeChain
	tail *NodeChain
}

type NodeChain struct {
	Value *model.Node
	Prev  *NodeChain
	Next  *NodeChain
}

func (nc *NodeChain) Remove() *NodeChain {
	var new *NodeChain
	if nc.Prev != nil {
		nc.Prev.Next = nc.Next
		new = nc.Prev
	} else {
		if nc.Next != nil {
			nc.Next.Prev = nil
		}
		new = nc.Next
	}

	nc = nil
	return new
}

func (nc *NodeChain) PushFront(n *model.Node) *NodeChain {
	new := &NodeChain{Value: n}
	new.Prev = nc
	new.Next = nc.Next
	if nc.Next != nil {
		nc.Next.Prev = new
	}
	nc.Next = new
	return new
}

type Session struct {
	//Store the nodes in the graph
	Graph []*model.Node
	//Representation of the density connected tree
	DCT []*model.Node
	//Store the cardinality of the partition when we explore the density connected tree from a starting node and excluding one edge
	//map[FromNode][ExcludeNode]Cardinality of the partition
	ExploreResults map[model.NodeID]map[model.NodeID]int
	//Store in order the node we are adding to the density conneted tree
	T []*model.Node
	//Contains the two node indexes of the edge with the smallest D-cut score
	minFrom, minTo model.NodeID
	//Smallest D-cut score
	minDcut float64
}

//DensityConnectedTree Create the density connected tree starting from a give node. If the first node is not provided, it randomly picks one.
func (s *Session) DensityConnectedTree(Graph []*model.Node, first *int) error {
	gSize := len(Graph)
	s.DCT = make([]*model.Node, gSize)
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
	tLength := 1

	for {
		//maxv = −1; p = null; q = null;
		maxv := float64(-1)
		var p, q *model.Node
		//while T.size < V.size do
		if tLength >= gSize {
			break
		}
		N := len(T)
		for i := 0; i < N; i++ {
			u := T[i]
			if len(u.Neighbors) == 0 {
				return fmt.Errorf("Node with index %s does not have any neighbors", u)
			}
			var availableNeighbor bool
			for _, vEdge := range u.Neighbors {
				if int(vEdge.To) >= len(Graph) {
					return fmt.Errorf("Node with index %d does not exist", vEdge.To)
				}
				//v = Γ(u).get(j);
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
					availableNeighbor = true
				}
			}
			if !availableNeighbor {
				T[i] = T[len(T)-1]
				T = T[:len(T)-1]
				N = len(T)
				i--
			}
		}
		p.Checked = true
		//After each iteration, we create a new edge in the Density Connected Tree.
		//Check is true, because we want to only check the Dcut bi-partition for one of the edge.

		if s.DCT[p.Index] == nil {
			s.DCT[p.Index] = &model.Node{
				Index: p.Index,
				Value: p.Value,
			}
		}
		s.DCT[p.Index].Neighbors = append(s.DCT[p.Index].Neighbors, &model.Edge{To: q.Index, Weight: maxv, Check: true})

		if s.DCT[q.Index] == nil {
			s.DCT[q.Index] = &model.Node{
				Index: q.Index,
				Value: q.Value,
			}
		}
		s.DCT[q.Index].Neighbors = append(s.DCT[q.Index].Neighbors, &model.Edge{To: p.Index, Weight: maxv})
		T = append(T, p)
		tLength++
	}
	s.Graph = Graph
	s.T = T
	return nil
}

//explore Returns the cardinality of the partition when we cut the edge between `node` and `exclude` in the density connected tree
func (s *Session) explore(node model.NodeID, exclude model.NodeID) int {
	val, ok := s.ExploreResults[node]

	if !ok {
		s.ExploreResults[node] = make(map[model.NodeID]int)
	} else {
		if storedScore, ok2 := val[exclude]; ok2 {
			return storedScore
		}
	}

	//-1 to exclude
	count := len(s.DCT[node].Neighbors) - 1
	for _, edge := range s.DCT[node].Neighbors {
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
func (s *Session) Dcut() (model.NodeID, model.NodeID, float64) {
	//For each nodeA we want to count the number of nodes in the partition if we break the edge between nodeA and nodeB
	//map[nodeA][nodeB]partitionSize
	s.ExploreResults = make(map[model.NodeID]map[model.NodeID]int)
	minDcut := math.Inf(1)
	minFrom := model.NodeID(-1)
	minTo := model.NodeID(-1)
	if len(s.DCT) > 1 {
		//For each edge in the Density Connected Tree, we evaluate the score of the two partitions defined after removing that edge
		dcut := float64(0)
		for index, node := range s.DCT {
			for _, e := range node.Neighbors {
				if e.Check {
					//Count partiion should also include the node itself --> + 1
					countParition := s.explore(model.NodeID(index), e.To) + 1
					//Dcut(C1, C2) = d(C1, C2)/min(|C1|, |C2|)
					dcut = e.Weight / float64(min(countParition, len(s.Graph)-countParition))

					if dcut < minDcut || (dcut == minDcut && ((model.NodeID(index) <= minFrom && e.To <= minTo) || (model.NodeID(index) <= minTo && e.To <= minFrom))) {
						minDcut = dcut
						minFrom = model.NodeID(index)
						minTo = e.To
					}
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
func (s *Session) extractParition(partition map[model.NodeID]model.NodeID, node, exclude, nextNodeId model.NodeID) (map[model.NodeID]model.NodeID, model.NodeID) {
	partition[node] = nextNodeId

	for _, edge := range s.DCT[node].Neighbors {
		if edge.To != exclude {
			nextNodeId++
			partition, nextNodeId = s.extractParition(partition, edge.To, node, nextNodeId)
		}
	}
	return partition, nextNodeId
}

//CreatePartition Returns a partition of a graph based on the results of a Dcut.
func (s *Session) CreatePartition(from, exclude model.NodeID) []*model.Node {
	paritionSize, ok := s.ExploreResults[from][exclude]
	//include from in the partition count
	paritionSize += 1
	if !ok {
		//If the count does not exsit, we can predict it
		paritionSize = len(s.Graph) - (s.ExploreResults[exclude][from] + 1)
	}

	log.Printf("Partition cardinality: %d (explore graph from vertex %s and ignore edge to %s)", paritionSize, s.Graph[from].Value, s.Graph[exclude].Value)

	partition1ID := make(map[model.NodeID]model.NodeID, paritionSize)
	partition1ID, _ = s.extractParition(partition1ID, from, exclude, 0)
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
	partition1 := s.CreatePartition(s.minFrom, s.minTo)
	partition2 := s.CreatePartition(s.minTo, s.minFrom)
	return partition1, partition2
}
