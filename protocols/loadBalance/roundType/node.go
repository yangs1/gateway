package roundType

import "fmt"

const DefaultCheckMaxErrNum = 2

type NodeInterface interface {
	OnHealthChange(bool)
	Get() Node
}

type Node struct {
	Ip              string
	Weight          int
	EffectiveWeight int
}

func (n *Node) OnHealthChange(checkRes bool) {
	if checkRes {
		if n.EffectiveWeight < DefaultCheckMaxErrNum {
			n.EffectiveWeight++
		}
	} else {
		if n.EffectiveWeight > 0 {
			n.EffectiveWeight--
		}
	}
}

type NodeHandler struct {
	nodes     []NodeInterface
	failNodes []NodeInterface
}

func (r *NodeHandler) Lists() []NodeInterface {
	return append(r.nodes, r.failNodes...)
}

func (r *NodeHandler) ConnWatcher() {
	var (
		nodes     []NodeInterface
		failNodes []NodeInterface
	)

	allNodes := r.Lists()
	for _, n := range allNodes {
		if n.Get().EffectiveWeight > 0 {
			nodes = append(nodes, n)
		} else {
			failNodes = append(failNodes, n)
		}
	}

	fmt.Println(fmt.Sprintf("success_node length: %d", len(nodes)))
	fmt.Println(fmt.Sprintf("fail_node length: %d", len(failNodes)))
	fmt.Println("================================")

	r.nodes = nodes
	r.failNodes = failNodes
}
