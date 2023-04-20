package protocols

//Server is comet server interface
type Server interface {
	Protocol() string
	Config() interface{}
	IsBind() bool
}
