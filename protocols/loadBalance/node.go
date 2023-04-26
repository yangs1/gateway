package loadBalance

import (
	"gateway/protocols/loadBalance/roundType"
	"net"
	"time"
)

const (
	DefaultCheckTimeout   = 3
	DefaultCheckMaxErrNum = 2
	DefaultCheckInterval  = 5
)

const (
	LbRandom int = iota
	LbRoundRobin
	LbRoundRobinWithWeight
	LbRandHash
)

type LoadBalance interface {
	Add(roundType.Node)
	Next(string) roundType.NodeInterface
	Lists() []roundType.NodeInterface
}

type LoadBalanceHandler struct {
	handler LoadBalance
}

func NewLoadBalance(lbType int) *LoadBalanceHandler {
	loadbalance := new(LoadBalanceHandler)

	switch lbType {
	case LbRandom:
		loadbalance.handler = new(roundType.RandomBalance)
	case LbRandHash:
		loadbalance.handler = new(roundType.HashBalance)
	case LbRoundRobin:
		loadbalance.handler = new(roundType.RoundRobinBalance)
	case LbRoundRobinWithWeight:
		loadbalance.handler = new(roundType.WeightRoundRobinBalance)
	}

	return loadbalance
}

func (handler *LoadBalanceHandler) Watch() {
	go func() {
		t := time.NewTicker(time.Second * DefaultCheckInterval)

		for {
			select {
			case <-t.C:
				handler.healthChecker()
			}
		}
	}()
}

func (loadbalance *LoadBalanceHandler) healthChecker() {
	nodes := loadbalance.handler.Lists()

	for _, n := range nodes {
		conn, err := net.DialTimeout("tcp", n.Get().Ip, time.Duration(DefaultCheckTimeout)*time.Second)
		//todo http statuscode
		if err == nil {
			conn.Close()
			n.OnHealthChange(true)
		} else {
			n.OnHealthChange(false)
		}
	}
}
