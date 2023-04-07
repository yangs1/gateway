package loadBalance

import (
	"gateway/global"
	"hash/crc32"
	"net/http"
)

type RandHashBalance struct {
	curIndex int
	rss      []string
	failRss  []string
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
	client := http.Client{}
	allRss := append(r.rss, r.failRss...)
	successRss := make([]string, 0)
	failRss := make([]string, 0)

	for _, target := range allRss {
		res, err := client.Head(target)
		if res != nil {
			defer res.Body.Close()
		}
		if err != nil {
			failRss = append(failRss, target)
			continue
		}
		if res.StatusCode >= 200 && res.StatusCode < 500 {
			successRss = append(successRss, target)
		} else {
			failRss = append(failRss, target)
		}
	}

	r.rss = successRss
	r.failRss = failRss

	global.Logger.Info("=========== success =============")
	for _, rs := range successRss {
		global.Logger.Info(rs)
	}
}
