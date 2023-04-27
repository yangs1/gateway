package roundType

type WeightRoundRobinBalance struct {
	curIndex int
	NodeHandler
}

type WeightRoundRobinNode struct {
	Node
	currentWeight int
}

func (r WeightRoundRobinNode) Get() Node {
	return r.Node
}

func (r *WeightRoundRobinBalance) Add(n Node) {
	r.nodes = append(r.nodes, &WeightRoundRobinNode{n, 5})
}

func (r *WeightRoundRobinBalance) Next(_ ...string) NodeInterface {
	nodes := r.nodes

	total := 0
	var best *WeightRoundRobinNode
	for i := 0; i < len(nodes); i++ {
		w := nodes[i].(*WeightRoundRobinNode)
		//step 1 统计所有有效权重之和
		total += w.EffectiveWeight

		//step 2 变更节点临时权重为的节点临时权重+节点有效权重
		w.currentWeight += w.EffectiveWeight

		//step 3 有效权重默认与权重相同，通讯异常时-1, 通讯成功+1，直到恢复到weight大小
		if w.EffectiveWeight < w.Weight {
			w.EffectiveWeight++
		}
		//step 4 选择最大临时权重点节点
		if best == nil || w.currentWeight > best.currentWeight {
			best = w
		}
	}
	if best == nil {
		return nil
	}
	//step 5 变更临时权重为 临时权重-有效权重之和
	best.currentWeight -= total
	return best
}
