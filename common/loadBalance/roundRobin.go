package loadBalance

type RoundRobinBalance struct {
	curIndex int
	rss      []string
}

func (r *RoundRobinBalance) Add(nodes ...BalanceNode) error {
	for _, n := range nodes {
		r.rss = append(r.rss, n.Addr)
	}

	return nil
}

func (r *RoundRobinBalance) Next() string {
	if len(r.rss) == 0 {
		return ""
	}
	lens := len(r.rss)
	if r.curIndex >= lens {
		r.curIndex = 0
	}
	curAddr := r.rss[r.curIndex]
	r.curIndex = (r.curIndex + 1) % lens
	return curAddr
}

func (r *RoundRobinBalance) Get(key string) (string, error) {
	return r.Next(), nil
}
