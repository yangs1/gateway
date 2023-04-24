package loadBalance

const (
	DefaultCheckTimeout   = 3
	DefaultCheckMaxErrNum = 2
	DefaultCheckInterval  = 5
)

type LoadBalanceNode interface {
	Add(...*Node) error
	Get(string) *Node
}

type Node struct {
	Ip     string
	Weight int
}
