package roundType

import (
	"gateway/protocols/loadBalance"
)

type RoundRobinBalance struct {
	curIndex int
	nodes    []*RoundRobinNode
}

type RoundRobinNode struct {
	*loadBalance.Node
	currentWeight uint
}

func (r *RoundRobinBalance) Add(nodes ...*loadBalance.Node) error {

	for _, n := range nodes {
		r.nodes = append(r.nodes, &RoundRobinNode{n, loadBalance.DefaultCheckTimeout})
	}

	return nil
}

func (r *RoundRobinBalance) Next() *RoundRobinNode {
	lens := len(r.nodes)

	if lens == 0 {
		return nil
	}
	// todo 需要 atomic 操作
	r.curIndex = (r.curIndex + 1) % lens

	return r.nodes[r.curIndex]
}

func (r *RoundRobinBalance) Get(string) *loadBalance.Node {
	if n := r.Next(); n != nil {
		return n.Node
	}

	return nil
}
