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
	Next(...string) roundType.NodeInterface
	Lists() []roundType.NodeInterface
}

type LoadBalanceHandler struct {
	Handler LoadBalance
}

func NewLoadBalance(lbType int) *LoadBalanceHandler {
	loadbalance := new(LoadBalanceHandler)

	switch lbType {
	case LbRandom:
		loadbalance.Handler = new(roundType.RandomBalance)
	case LbRandHash:
		loadbalance.Handler = new(roundType.HashBalance)
	case LbRoundRobin:
		loadbalance.Handler = new(roundType.RoundRobinBalance)
	case LbRoundRobinWithWeight:
		loadbalance.Handler = new(roundType.WeightRoundRobinBalance)
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
	nodes := loadbalance.Handler.Lists()

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
