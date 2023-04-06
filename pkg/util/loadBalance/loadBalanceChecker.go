package loadBalance

import (
	"time"
)

type LoadBalanceChecker struct {
	Lb               LoadBalance
	effective_weight int
}

func NewLbChecker(lb LoadBalance, weight int) *LoadBalanceChecker {
	return &LoadBalanceChecker{
		Lb:               lb,
		effective_weight: weight,
	}
}

func (checker *LoadBalanceChecker) HttpWatch() {
	go func() {
		t := time.NewTicker(time.Second * 5)

		for {
			select {
			case <-t.C:
				checker.Lb.Check()
			}
		}
	}()
}
