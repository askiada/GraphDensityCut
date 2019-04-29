package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"sync"

	metric "github.com/askiada/GraphDensityCut/src/distance"

	"github.com/askiada/GraphDensityCut/src/model"
)

type Session struct {
	G              []*model.Node
	DCTEdges       map[int][]*model.Edge
	ExploreResults map[int]map[int]int
	//ExploreResults map[int]map[int]*PartitionInfo
	T              []*model.Node
	minFrom, minTo int
	minDcut        float64
	id             int
}

type Candidate struct {
	p, q *model.Node
	sim  float64
}

func (s *Session) DensityConnectedTree(G []*model.Node, first *int) error {
	//s.DCTCount = make(map[int]int)
	s.DCTEdges = make(map[int][]*model.Edge)
	gSize := len(G)
	//T = null;
	var T []*model.Node
	//Set ∀v ∈ V as unchecked (v.checked = false); --> Zero value for boolean
	//Randomly selected one node u ∈ V ;
	if first == nil {
		tmp := rand.Intn(len(G))
		first = &tmp
	}
	//Set u.checked = true;
	G[*first].Checked = true

	//u.connect = null, and u.density = null; --> Zero value for pointer connect. Density is equal to 0 by default, and it does not matter

	//T.insert(u);
	metric.Init(len(G))
	T = append(T, G[*first])
	for true {
		//maxv = −1; p = null; q = null;
		maxv := float64(-1)
		var p, q *model.Node
		//while T.size < V.size do
		if len(T) >= gSize {
			break
		}
		wg := &sync.WaitGroup{}
		rChan := make(chan *Candidate)
		go func() {
			for {
				if r, open := <-rChan; open {
					if r.sim > maxv {
						maxv = r.sim
						p = r.p
						q = r.q
					}
				} else {
					break
				}
			}
		}()
		//for j = 1 to Γ(u).size do
		for i := range T {
			wg.Add(1)
			u := T[i]
			if len(u.Neighbors) == 0 {
				return fmt.Errorf("Node with index %s does not have any neighbors", u)
			}
			go func(wg2 *sync.WaitGroup, g2 []*model.Node, u2 *model.Node) {
				defer wg.Done()
				wg.Add(len(u.Neighbors))
				for j := range u2.Neighbors {
					//v = Γ(u).get(j);
					vEdge := u2.Neighbors[j]
					if vEdge.To >= len(g2) {
						panic("Node with index does not exist" + strconv.Itoa(vEdge.To))
					}
					go func(wg1 *sync.WaitGroup, vE *model.Edge, g1 []*model.Node, u1 *model.Node) {
						defer wg1.Done()

						v := g1[vE.To]
						//if v.checked == false then
						if !v.Checked {
							//If we have already computed the NodeSimilarity for an edge, we can use the score from the previous computation
							if vE.NodeSimilarity == nil {
								tmp := metric.NodeSim(u1, v, vE.Weight)
								vE.NodeSimilarity = &tmp
							}
							//if s(u, v) > maxv then
							if *vE.NodeSimilarity > maxv {
								maxv = *vE.NodeSimilarity
								p = v
								q = u1
							}

							c := &Candidate{p: v, q: u1, sim: *vE.NodeSimilarity}
							rChan <- c

						}
					}(wg2, vEdge, g2, u2)
				}
			}(wg, G, u)
		}
		wg.Wait()
		close(rChan)
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
	s.G = G
	s.T = T
	return nil
}

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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

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

func (s *Session) CreatePartition(from, exclude int) []*model.Node {
	paritionSize, ok := s.ExploreResults[from][exclude]
	//include from in the partition count
	paritionSize += 1
	if !ok {
		//If the count does not exsit, we can predict it
		paritionSize = len(s.G) - (s.ExploreResults[exclude][from] + 1)
	}
	partition1ID := make(map[int]int, paritionSize)
	partition1ID = s.extractParition(partition1ID, from, exclude)
	partition1 := make([]*model.Node, paritionSize)

	for node, idx := range partition1ID {
		partition1[idx] = &model.Node{}
		*partition1[idx] = *s.G[node]
		partition1[idx].Checked = false
		partition1[idx].Index = idx
		partition1[idx].Neighbors = nil
		newEdges := []*model.Edge{}
		for _, e := range s.G[node].Neighbors {
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
	partition1 := s.CreatePartition(s.minFrom, s.minTo)
	//Reset index of partition
	s.id = 0
	partition2 := s.CreatePartition(s.minTo, s.minFrom)
	return partition1, partition2
}

func main() {
	//Nothing for now
}
