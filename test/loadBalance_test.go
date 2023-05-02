package main

import (
	"fmt"
	"gateway/protocols/loadBalance"
	"gateway/protocols/loadBalance/roundType"
	"testing"
	"time"
)

func TestLb(t *testing.T) {
	lb := loadBalance.NewLoadBalance(0)
	lb.Handler.Add(roundType.Node{
		Ip:              "127.0.0.1:9092",
		Weight:          2,
		EffectiveWeight: 2,
	})
	lb.Handler.Add(roundType.Node{
		Ip:              "127.0.0.1:9091",
		Weight:          0,
		EffectiveWeight: 2,
	})

	c := lb.Handler.Lists()
	for _, n := range c {
		n.OnHealthChange(false)
	}
	fmt.Println(lb.Handler.Next().Get())

}

func TestLbWatch(t *testing.T) {
	lb := loadBalance.NewLoadBalance(0)
	lb.Handler.Add(roundType.Node{
		Ip:              "127.0.0.1:9092",
		Weight:          2,
		EffectiveWeight: 2,
	})
	lb.Handler.Add(roundType.Node{
		Ip:              "127.0.0.1:9091",
		Weight:          0,
		EffectiveWeight: 2,
	})
	lb.Handler.Add(roundType.Node{
		Ip:              "127.0.0.1:9093",
		Weight:          3,
		EffectiveWeight: 2,
	})
	lb.Watch()
	for {
		time.Sleep(5 * time.Second)
	}
}
