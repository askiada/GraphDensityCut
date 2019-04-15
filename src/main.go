package main

import (
	"math"
	"math/rand"

	metric "github.com/askiada/Dcut/src/distance"

	"github.com/askiada/Dcut/src/model"
)

type Session struct {
	DCTCount       map[int]int
	DCTEdges       map[int][]*model.Edge
	ExploreResults map[int]map[int]int
}

func (s *Session) DensityConnectedTree(g []*model.Node, first *int) []*model.Node {

	s.DCTCount = make(map[int]int)
	s.DCTEdges = make(map[int][]*model.Edge)
	gSize := len(g)
	//T = null;
	var T []*model.Node
	//Set ∀v ∈ V as unchecked (v.checked = f alse); --> Zero value for boolean
	//Randomly selected one node u ∈ V ;
	if first == nil {
		*first = rand.Intn(len(g))
	}
	//Set u.checked = true;
	g[*first].Checked = true

	//u.connect = null, and u.density = null; --> Zero value for pointers

	//T.insert(u);

	T = append(T, g[*first])
	for true {
		maxv := float64(-1)
		var p, q *model.Node

		//while T.size < V.size do
		if len(T) >= gSize {
			break
		}

		for i := range T {
			u := T[i]
			for j := range u.Neighbors {
				v := u.Neighbors[j]
				if !g[v.To].Checked {
					suv := metric.JaccardCoeff(u, g[v.To])
					if suv > maxv {
						maxv = suv
						p = g[v.To]
						q = u
					}
				}
			}
		}

		p.Checked = true
		p.Connect = q
		p.Density = maxv

		s.DCTCount[p.Value]++
		s.DCTCount[q.Value]++
		s.DCTEdges[p.Value] = append(s.DCTEdges[p.Value], &model.Edge{To: q.Value, Weight: maxv, Check: true})
		s.DCTEdges[q.Value] = append(s.DCTEdges[q.Value], &model.Edge{To: p.Value, Weight: maxv})
		T = append(T, p)
	}
	return T
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
	count := s.DCTCount[node] - 1
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
	s.ExploreResults = make(map[int]map[int]int)
	minDcut := math.Inf(1)
	minFrom := -1
	minTo := -1

	dcut := float64(0)
	for node, edges := range s.DCTEdges {
		for _, e := range edges {
			if e.Check {
				countParition := s.explore(node, e.To) + 1
				dcut = e.Weight / float64(min(countParition, len(s.DCTCount)-countParition))

				if dcut < minDcut && dcut != 0 {
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

}
