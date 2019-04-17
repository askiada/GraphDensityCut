package main

import (
	"fmt"
	"math"
	"math/rand"

	metric "github.com/askiada/GraphDensityCut/src/distance"

	"github.com/askiada/GraphDensityCut/src/model"
)

type Session struct {
	DCTEdges       map[int][]*model.Edge
	ExploreResults map[int]map[int]int
	T              []*model.Node
}

func (s *Session) DensityConnectedTree(g []*model.Node, first *int) error {
	//s.DCTCount = make(map[int]int)
	s.DCTEdges = make(map[int][]*model.Edge)
	gSize := len(g)
	//T = null;
	var T []*model.Node
	//Set ∀v ∈ V as unchecked (v.checked = false); --> Zero value for boolean
	//Randomly selected one node u ∈ V ;
	if first == nil {
		tmp := rand.Intn(len(g))
		first = &tmp
	}
	//Set u.checked = true;
	g[*first].Checked = true

	//u.connect = null, and u.density = null; --> Zero value for pointer connect. Density is equal to 0 by default, and it does not matter

	//T.insert(u);
	metric.Init(len(g))
	T = append(T, g[*first])

	for true {
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
				if vEdge.To >= len(g) {
					return fmt.Errorf("Node with index %d does not exist", vEdge.To)
				}
				v := g[vEdge.To]
				//if v.checked == false then
				if !v.Checked {
					//If we have already computed the NodeSimilarity for an edge, we can use the score from the previous computation
					if vEdge.NodeSimilarity == nil {

						tmp := metric.NodeSim(u, v, vEdge.Weight)
						vEdge.NodeSimilarity = &tmp
						fmt.Println(u.Value+1, v.Value+1, tmp)
					}
					//if s(u, v) > maxv then
					if *vEdge.NodeSimilarity > maxv {
						maxv = *vEdge.NodeSimilarity
						p = v
						q = u
					}
					//fmt.Println(maxv)
				}
			}
		}
		p.Checked = true
		p.Connect = q
		p.Density = maxv
		//After each iteration, we create a new edge in the Density Connected Tree.
		//Check is true, because we want to only check the Dcut bi-partition for one of the edge.
		s.DCTEdges[p.Value] = append(s.DCTEdges[p.Value], &model.Edge{To: q.Value, Weight: maxv, Check: true})
		s.DCTEdges[q.Value] = append(s.DCTEdges[q.Value], &model.Edge{To: p.Value, Weight: maxv})
		//T.insert(p);
		T = append(T, p)
	}
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

				if dcut < minDcut {
					minDcut = dcut
					minFrom = node
					minTo = e.To
				}
			}
		}
	}
	return minFrom, minTo, minDcut
}

func main() {
	//Nothing for now
}
