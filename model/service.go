package model

type ServiceDetail struct {
	Addr     string
	Weight   int
	LoadType int
	isHttps  bool
	isTcp    bool
	Timeout  int
}
