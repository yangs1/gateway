package loadBalance

import (
	"fmt"
	"gateway/model/loadBalance"
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
	ConnWatcher()
}

type LoadBalanceHandler struct {
	roundType int
	Handler   LoadBalance
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

	loadbalance.roundType = lbType

	return loadbalance
}

func (handler *LoadBalanceHandler) AddNode(node loadBalance.ServerHttp) {
	handler.Handler.Add(roundType.Node{
		Ip:              node.Ip,
		Weight:          node.Weight,
		EffectiveWeight: DefaultCheckMaxErrNum,
	})
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
	nodes := handler.Handler.Lists()

	for _, n := range nodes {
		conn, err := net.DialTimeout("tcp", n.Get().Ip, time.Duration(DefaultCheckTimeout)*time.Second)
		//todo http statuscode
		// todo 检测失败最大次数
		if err == nil {
			conn.Close()
			n.OnHealthChange(true)

			fmt.Println(fmt.Sprintf("node: %s is success", n.Get().Ip))

		} else {
			fmt.Println(fmt.Sprintf("node: %s is fail", n.Get().Ip))

			n.OnHealthChange(false)
		}
	}
	handler.Handler.ConnWatcher()
}
