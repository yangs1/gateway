package protocols

import (
	"fmt"
	"sync"
)

type FactoryFun func() Server

type serverIoc struct {
	provider map[string]FactoryFun
	// instance 存储具体的实例
	instances map[string]Server
	// lock 用于锁住对容器的变更操作
	lock sync.RWMutex
}

var ServerIoc *serverIoc

// 设置简易 ioc 容器
func init() {
	ServerIoc = &serverIoc{
		provider:  make(map[string]FactoryFun),
		instances: make(map[string]Server),
		lock:      sync.RWMutex{},
	}
}

// Register registers a new output type.
func (ioc *serverIoc) Register(name string, f FactoryFun) {
	ioc.lock.Lock()
	defer ioc.lock.Unlock()

	if ioc.provider[name] != nil || ioc.instances[name] != nil {
		panic(fmt.Errorf("server type  '%v' exists already", name))
	}
	ioc.provider[name] = f
}

// Bind registers a new output type.
func (ioc *serverIoc) Bind(name string, s Server) {
	ioc.lock.Lock()
	defer ioc.lock.Unlock()

	if ioc.provider[name] != nil || ioc.instances[name] != nil {
		panic(fmt.Errorf("server type  '%v' exists already", name))
	}

	ioc.instances[name] = s
}

// Load get a func or instance from ioc.
func (ioc *serverIoc) Load(name string) (Server, error) {
	ioc.lock.Lock()
	defer ioc.lock.Unlock()

	instance := ioc.instances[name]
	if instance != nil {
		return instance, nil
	}

	factory := ioc.provider[name]
	if factory != nil {
		instance := factory()
		if instance.IsBind() {
			ioc.instances[name] = instance
		}
		return instance, nil
	}

	return nil, fmt.Errorf("server type %v undefined", name)
}
