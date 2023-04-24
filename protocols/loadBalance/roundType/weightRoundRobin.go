package roundType

import "gateway/protocols/loadBalance"

type WeightRoundRobinBalance struct {
	curIndex int
	nodes    []*WeightRoundRobinNode
}

type WeightRoundRobinNode struct {
	*loadBalance.Node
	currentWeight   int
	effectiveWeight int //有效权重

}

func (r *WeightRoundRobinBalance) Add(nodes ...*loadBalance.Node) error {

	for _, n := range nodes {
		r.nodes = append(r.nodes, &WeightRoundRobinNode{n, n.Weight, n.Weight})
	}

	return nil
}

func (r *WeightRoundRobinBalance) Next() *WeightRoundRobinNode {
	total := 0
	var best *WeightRoundRobinNode
	for i := 0; i < len(r.nodes); i++ {
		w := r.nodes[i]
		//step 1 统计所有有效权重之和
		total += w.effectiveWeight

		//step 2 变更节点临时权重为的节点临时权重+节点有效权重
		w.currentWeight += w.effectiveWeight

		//step 3 有效权重默认与权重相同，通讯异常时-1, 通讯成功+1，直到恢复到weight大小
		if w.effectiveWeight < w.Weight {
			w.effectiveWeight++
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

func (r *WeightRoundRobinBalance) Get(_ string) *loadBalance.Node {
	if n := r.Next(); n != nil {
		return n.Node
	}

	return nil
}
