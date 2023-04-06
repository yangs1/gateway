package loadBalance

import (
	"hash/crc32"
)

type RandHashBalance struct {
	curIndex int
	rss      []string
}

func (r *RandHashBalance) Add(nodes ...BalanceNode) error {

	for _, n := range nodes {
		r.rss = append(r.rss, n.Addr)
	}

	return nil
}

func (r *RandHashBalance) Get(key string) (string, error) {

	index := int(crc32.ChecksumIEEE([]byte(key))) % len(r.rss)

	return r.rss[index], nil
}

func (r *RandHashBalance) Check() {

}
