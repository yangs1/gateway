package initialize

import (
	"gateway/global"
	"gateway/pkg/util"
	"gateway/pkg/util/loadBalance"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
	"net"
	"net/http"
	"sync"
	"time"
)

var ServiceManagerHandler *ServiceManager

//专门用于处理service的对象
type ServiceManager struct {
	loadBalanceHandler loadBalance.LoadBalance
	transportHandler   *http.Transport
	pxyCfg             *ProxyConfig
	locker             sync.RWMutex //加个锁
	init               sync.Once    //一次性加载
	err                error
}

type ProxyConfig struct {
	LoadType              int                       `yaml:"load_type"`
	EffectiveWeight       int                       `yaml:"effective_weight"`
	MaxIdleConns          int                       `yaml:"max_idle"`
	IdleConnTimeout       int                       `yaml:"idle_timeout"`
	ConnectTimeout        int                       `yaml:"connect_timeout"`
	ResponseHeaderTimeout int                       `yaml:"header_timeout"`
	rss                   []loadBalance.BalanceNode `yaml:"proxy_nodes"`
}

func init() {
	ServiceManagerHandler = &ServiceManager{
		init: sync.Once{},
	}
}

func InitLoadBalance() error {
	ServiceManagerHandler.init.Do(func() {
		vcfg, err := util.ReadConfigFile("proxy")
		if err != nil {
			global.Logger.Error("loadbalance 配置文件读取失败", zap.Error(err))
			ServiceManagerHandler.err = err
			return
		}
		var proxyConfig ProxyConfig
		err = vcfg.Unmarshal(&proxyConfig, func(config *mapstructure.DecoderConfig) {
			config.TagName = "yaml"
		})

		if err != nil {
			global.Logger.Error("加载配置失败", zap.Error(err))
			ServiceManagerHandler.err = err
			return
		}

		// TODO Unmarshal 没解析出来，之后看看什么问题
		if err := mapstructure.Decode(vcfg.Get("proxy_nodes"), &proxyConfig.rss); err != nil {
			global.Logger.Error("加载代理配置失败", zap.Error(err))
			ServiceManagerHandler.err = err
			return
		}

		ServiceManagerHandler.pxyCfg = &proxyConfig

	})

	return ServiceManagerHandler.err
}

// 获取负载均衡器
func (manager *ServiceManager) GetLoadBalancer() loadBalance.LoadBalance {
	if manager.loadBalanceHandler != nil {
		return manager.loadBalanceHandler
	}

	manager.locker.Lock()
	defer manager.locker.Unlock()

	manager.loadBalanceHandler = loadBalance.NewLoadBalance(manager.pxyCfg.LoadType)
	manager.loadBalanceHandler.Add(manager.pxyCfg.rss...)

	httpChecker := loadBalance.NewLbChecker(manager.loadBalanceHandler, 5)
	httpChecker.HttpWatch()
	
	return manager.loadBalanceHandler
}

// 获取连接池
//@see https://www.bbsmax.com/A/gGdX3GMG54/
func (manager *ServiceManager) GetTrans() *http.Transport {
	if manager.transportHandler != nil {
		return manager.transportHandler
	}

	manager.locker.Lock()
	defer manager.locker.Unlock()

	pxyCfg := manager.pxyCfg
	manager.transportHandler = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   time.Duration(pxyCfg.ConnectTimeout) * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          pxyCfg.MaxIdleConns,
		IdleConnTimeout:       time.Duration(pxyCfg.IdleConnTimeout) * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: time.Duration(pxyCfg.ResponseHeaderTimeout) * time.Second,
	}

	return manager.transportHandler
}
