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
	currentWeight uint
}

func (r *HashBalance) Add(nodes ...*loadBalance.Node) error {

	for _, n := range nodes {
		r.nodes = append(r.nodes, &HashNode{n, loadBalance.DefaultCheckTimeout})
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
