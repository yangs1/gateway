package loadBalance

type RoundRobinWithWeightBalance struct {
	curIndex int
	rss      []*WeightNode
}

type WeightNode struct {
	BalanceNode
	curWeight       int //节点当前权重
	effectiveWeight int //有效权重
}

func (r *RoundRobinWithWeightBalance) Add(nodes ...BalanceNode) error {

	for _, n := range nodes {
		node := &WeightNode{BalanceNode: n}
		node.effectiveWeight = node.Weight
		r.rss = append(r.rss, node)
	}

	return nil
}

func (r *RoundRobinWithWeightBalance) Next() string {
	total := 0
	var best *WeightNode
	for i := 0; i < len(r.rss); i++ {
		w := r.rss[i]
		//step 1 统计所有有效权重之和
		total += w.effectiveWeight

		//step 2 变更节点临时权重为的节点临时权重+节点有效权重
		w.curWeight += w.effectiveWeight

		//step 3 有效权重默认与权重相同，通讯异常时-1, 通讯成功+1，直到恢复到weight大小
		if w.effectiveWeight < w.Weight {
			w.effectiveWeight++
		}
		//step 4 选择最大临时权重点节点
		if best == nil || w.curWeight > best.curWeight {
			best = w
		}
	}
	if best == nil {
		return ""
	}
	//step 5 变更临时权重为 临时权重-有效权重之和
	best.curWeight -= total
	return best.Addr
}

func (r *RoundRobinWithWeightBalance) Get(key string) (string, error) {
	return r.Next(), nil
}

func (r *RoundRobinWithWeightBalance) Check() {

}
