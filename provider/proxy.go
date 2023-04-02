package provider

import (
	"hash/crc32"
	"math/rand"
	"sort"
	"time"
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

func initProxy() *ProxyManager {
	manager := &ProxyManager{}

	server1 := &HttpServer{
		Host:   "127.0.0.1:9090",
		Weight: 4,
	}

	server2 := &HttpServer{
		Host:   "127.0.0.1:9091",
		Weight: 2,
	}

	manager.Servers = append(manager.Servers, server1)
	manager.Servers = append(manager.Servers, server2)

	for _, s := range manager.Servers {
		manager.TotalWeight += s.Weight
	}

	return manager
}

// 随机算法
func (manager *ProxyManager) SelectByRand() *HttpServer {
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(manager.Servers))

	return manager.Servers[index]
}

// ip_hash 算法
func (manager *ProxyManager) SelectByHash(ip string) *HttpServer {
	index := int(crc32.ChecksumIEEE([]byte(ip))) % len(manager.Servers)

	return manager.Servers[index]
}

// rand with weight 加权算法
func (manager *ProxyManager) SelectByWeightRand() *HttpServer {

	total := 0
	totalGrads := make([]int, len(manager.Servers))

	for _, s := range manager.Servers {
		total += s.Weight
		totalGrads = append(totalGrads, total)
	}

	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(total)

	for index, num := range totalGrads {
		if num > randNum {
			return manager.Servers[index]
		}
	}

	return manager.Servers[0]
}

// 轮询算法
func (manager *ProxyManager) SelectByLoop() *HttpServer {
	server := manager.Servers[manager.CurIndex]

	manager.CurIndex = (manager.CurIndex + 1) % len(manager.Servers)

	return server
}

// 加权轮询
func (manager *ProxyManager) SelectByLoopWeight() *HttpServer {
	server := manager.Servers[0]
	sum := 0

	for i, s := range manager.Servers {
		sum += s.Weight
		if manager.CurIndex < sum {
			server = manager.Servers[i]
			if manager.CurIndex == sum-1 && i != len(manager.Servers)-1 {
				manager.CurIndex++
			} else {
				manager.CurIndex = (manager.CurIndex + 1) % sum
			}
			break
		}
	}

	return server
}

// 平滑轮询算法
func (manager *ProxyManager) SelectByLoopLevelWeight() *HttpServer {

	for _, s := range manager.Servers {
		s.CurWeight += s.Weight + s.CurWeight
	}
	sort.Sort(manager.Servers)

	server := manager.Servers[0]
	server.CurWeight = server.CurWeight - manager.TotalWeight

	return server
}
