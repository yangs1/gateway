package loadBalance

type node struct {
	Ip              string
	Weight          uint
	EffectiveWeight int
}
