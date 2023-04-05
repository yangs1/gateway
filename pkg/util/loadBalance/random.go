package loadBalance

import (
	"math/rand"
	"time"
)

type RandomBalance struct {
	curIndex int
	rss      []string
}

func (r *RandomBalance) Add(nodes ...BalanceNode) error {

	for _, n := range nodes {
		r.rss = append(r.rss, n.Addr)
	}

	return nil
}

func (r *RandomBalance) Next() string {
	if len(r.rss) == 0 {
		return ""
	}

	rand.Seed(time.Now().UnixNano())
	r.curIndex = rand.Intn(len(r.rss))

	return r.rss[r.curIndex]
}

func (r *RandomBalance) Get(key string) (string, error) {
	return r.Next(), nil
}
