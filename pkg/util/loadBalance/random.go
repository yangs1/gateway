package loadBalance

import (
	"gateway/global"
	"math/rand"
	"net/http"
	"time"
)

type RandomBalance struct {
	curIndex int
	rss      []string
	failRss  []string
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

func (r *RandomBalance) Check() {
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

	global.Logger.Info("=====================Success=================================================")
	for _, tar := range successRss {
		global.Logger.Info(tar)
	}
	global.Logger.Info("=========================Fail=============================================")
	for _, tar := range failRss {
		global.Logger.Info(tar)
	}

}
