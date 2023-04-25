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

type LoadBalanceHandler struct {
	nodes BalanceNode
}

type BalanceNode interface {
	Add(...Node) error
	Get(string) *Node
	Nodes() []*Node
}

type Node struct {
	Ip              string
	Weight          int
	EffectiveWeight int
}

func (n *Node) OnHealthChange(checkRes bool) {
	if checkRes {
		if n.EffectiveWeight < DefaultCheckMaxErrNum {
			n.EffectiveWeight++
		}
	} else {
		if n.EffectiveWeight > 0 {
			n.EffectiveWeight--
		}
	}
}

func NewLoadBalance(lbType int) *LoadBalanceHandler {
	var lb BalanceNode

	switch lbType {
	case LbRandom:
		lb = &roundType.RandomBalance{}
	case LbRoundRobin:
		lb = &roundType.RoundRobinBalance{}
	case LbRoundRobinWithWeight:
		lb = &roundType.WeightRoundRobinBalance{}
	case LbRandHash:
		lb = &roundType.HashBalance{}
	}

	return &LoadBalanceHandler{
		nodes: lb,
	}
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

func (handler *LoadBalanceHandler) healthChecker() {
	nodes := handler.nodes.Nodes()

	for _, n := range nodes {
		conn, err := net.DialTimeout("tcp", n.Ip, time.Duration(DefaultCheckTimeout)*time.Second)
		//todo http statuscode
		if err == nil {
			conn.Close()
			n.OnHealthChange(true)
		} else {
			n.OnHealthChange(false)
		}
	}
}
