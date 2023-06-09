package roundType

import (
	"hash/crc32"
)

type HashBalance struct {
	NodeHandler
}

type HashNode struct {
	Node
}

func (r HashNode) Get() Node {
	return r.Node
}

func (r *HashBalance) Add(n Node) {
	r.nodes = append(r.nodes, &HashNode{n})
}

func (r *HashBalance) Next(key ...string) NodeInterface {
	ip := ""
	if len(key) > 0 {
		ip = key[0]
	}
	nodes := r.nodes

	index := int(crc32.ChecksumIEEE([]byte(ip))) % len(nodes)

	if n := nodes[index]; n != nil {
		return n
	}

	return nil
}
