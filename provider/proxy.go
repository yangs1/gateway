package provider

import (
	"gateway/common/loadBalance"
)

type ProxyServers []*HttpServer

type HttpServer struct {
	Host      string
	Weight    int
	CurWeight int
}

type ProxyManager struct {
	Servers     ProxyServers
	CurIndex    int
	TotalWeight int
}

func (servers ProxyServers) Len() int { return len(servers) }
func (servers ProxyServers) Swap(i, j int) {
	servers[i], servers[j] = servers[j], servers[i]
}
func (servers ProxyServers) Less(i, j int) bool {
	return servers[i].Weight > servers[j].Weight
}

func initProxy() loadBalance.LoadBalance {

	lbManager := loadBalance.NewLoadBalance(loadBalance.LbRoundRobinWithWeiht)

	lbManager.Add(loadBalance.BalanceNode{
		Addr:   "127.0.0.1:9090",
		Weight: 2,
	},
		loadBalance.BalanceNode{
			Addr:   "127.0.0.1:9091",
			Weight: 5,
		})

	return lbManager
}
