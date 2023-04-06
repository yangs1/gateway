package loadBalance

import "net/http"

type RoundRobinBalance struct {
	curIndex int
	rss      []string
	failRss  []string
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

func (r *RoundRobinBalance) Check() {
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
		if res.StatusCode >= 200 && res.StatusCode < 400 {
			successRss = append(successRss, target)
		} else {
			failRss = append(failRss, target)
		}
	}

	r.rss = successRss
	r.failRss = failRss
}
