package roundType

import (
	"math/rand"
	"time"
)

type RandomBalance struct {
	NodeHandler
}

type RandomNode struct {
	Node
}

func (r RandomNode) Get() Node {
	return r.Node
}

func (r *RandomBalance) Add(n Node) {
	r.nodes = append(r.nodes, &RandomNode{n})
}

func (r *RandomBalance) Next(_ ...string) NodeInterface {
	nodes := r.nodes
	if len(nodes) == 0 {
		return nil
	}

	rand.Seed(time.Now().UnixNano())
	curIndex := rand.Intn(len(nodes))

	return nodes[curIndex]
}
