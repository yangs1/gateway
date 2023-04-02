package loadBalance

type LoadBalance interface {
	Add(...BalanceNode) error
	Get(string) (string, error)
}

// 节点数据
type BalanceNode struct {
	Addr   string
	Weight int
}

const (
	LbRandom int = iota
	LbRoundRobin
	LbRoundRobinWithWeiht
	LbRandHash
)

func NewLoadBalance(lbType int) LoadBalance {
	var lb LoadBalance

	switch lbType {
	case LbRandom:
		lb = &RandomBalance{}
	case LbRoundRobin:
		lb = &RoundRobinBalance{}
	case LbRoundRobinWithWeiht:
		lb = &RoundRobinWithWeightBalance{}
	case LbRandHash:
		lb = &RandHashBalance{}
	}

	return lb
}
