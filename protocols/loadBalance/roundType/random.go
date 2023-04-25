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
}

func (r *RandomBalance) Add(nodes ...loadBalance.Node) error {

	for _, n := range nodes {
		r.nodes = append(r.nodes, &RandomNode{&loadBalance.Node{Ip: n.Ip, Weight: n.Weight, EffectiveWeight: loadBalance.DefaultCheckMaxErrNum}})
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

func (r *RandomBalance) Nodes() (nodes []*loadBalance.Node) {
	for _, n := range r.nodes {
		nodes = append(nodes, n.Node)
	}

	return nodes
}
