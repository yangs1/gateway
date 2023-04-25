package roundType

import (
	"gateway/protocols/loadBalance"
	"hash/crc32"
)

type HashBalance struct {
	curIndex int
	nodes    []*HashNode
}

type HashNode struct {
	*loadBalance.Node
}

func (r *HashBalance) Add(nodes ...loadBalance.Node) error {

	for _, n := range nodes {
		r.nodes = append(r.nodes, &HashNode{&loadBalance.Node{Ip: n.Ip, Weight: n.Weight, EffectiveWeight: loadBalance.DefaultCheckMaxErrNum}})
	}

	return nil
}

func (r *HashBalance) Get(ip string) *loadBalance.Node {
	index := int(crc32.ChecksumIEEE([]byte(ip))) % len(r.nodes)

	if n := r.nodes[index]; n != nil {
		return n.Node
	}

	return nil
}

func (r *HashBalance) Nodes() (nodes []*loadBalance.Node) {
	for _, n := range r.nodes {
		nodes = append(nodes, n.Node)
	}

	return nodes
}
