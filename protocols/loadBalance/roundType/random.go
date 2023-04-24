package roundType

import (
	"gateway/protocols/loadBalance"
	"math/rand"
	"time"
)

type RandomBalance struct {
	curIndex int
	nodes    []*RandomNode
}

type RandomNode struct {
	*loadBalance.Node
	currentWeight int
}

func (r *RandomBalance) Add(nodes ...*loadBalance.Node) error {

	for _, n := range nodes {
		r.nodes = append(r.nodes, &RandomNode{n, loadBalance.DefaultCheckTimeout})
	}

	return nil
}

func (r *RandomBalance) Next() *RandomNode {
	if len(r.nodes) == 0 {
		return nil
	}

	rand.Seed(time.Now().UnixNano())
	r.curIndex = rand.Intn(len(r.nodes))

	return r.nodes[r.curIndex]
}

func (r *RandomBalance) Get(string) *loadBalance.Node {
	if n := r.Next(); n != nil {
		return n.Node
	}

	return nil
}
