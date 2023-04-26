package roundType

type RoundRobinBalance struct {
	curIndex int
	NodeHandler
}

type RoundRobinNode struct {
	Node
}

func (r RoundRobinNode) Get() Node {
	return r.Node
}

func (r *RoundRobinBalance) Add(n Node) {
	r.nodes = append(r.nodes, &RoundRobinNode{n})
}

func (r *RoundRobinBalance) Next(_ string) NodeInterface {
	nodes := r.nodes

	lens := len(nodes)

	if lens == 0 {
		return nil
	}
	// todo 需要 atomic 操作
	r.curIndex = (r.curIndex + 1) % lens

	return nodes[r.curIndex]
}
