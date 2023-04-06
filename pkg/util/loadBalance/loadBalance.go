package loadBalance

type LoadBalance interface {
	Add(...BalanceNode) error
	Get(string) (string, error)
	Check()
}

// 节点数据
type BalanceNode struct {
	Addr   string `yaml:"addr"`
	Weight int    `yaml:"weight"`
}

const (
	LbRandom int = iota
	LbRoundRobin
	LbRoundRobinWithWeight
	LbRandHash
)

func NewLoadBalance(lbType int) LoadBalance {
	var lb LoadBalance

	switch lbType {
	case LbRandom:
		lb = &RandomBalance{}
	case LbRoundRobin:
		lb = &RoundRobinBalance{}
	case LbRoundRobinWithWeight:
		lb = &RoundRobinWithWeightBalance{}
	case LbRandHash:
		lb = &RandHashBalance{}
	}

	return lb
}
