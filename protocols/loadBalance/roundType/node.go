package roundType

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
