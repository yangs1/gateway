package initialize

import (
	"gateway/common/loadBalance"
	"gateway/global"
	"gateway/pkg/util"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
	"sync"
)

//专门用于处理service的对象
type ServiceManager struct {
	lb   map[string]loadBalance.LoadBalance
	init sync.Once //一次性加载
	err  error
}

func InitLoadBalance() error {
	global.ServiceManagerHandler.init.Do(func() {
		vcfg, err := util.ReadConfigFile("proxy")
		if err != nil {
			global.Logger.Error("general配置文件读取失败", zap.Error(err))
			global.ServiceManagerHandler.err = err
			return
		}

		var loadBalanceNodes []loadBalance.BalanceNode

		if err := mapstructure.Decode(vcfg.Get("http_proxy"), loadBalanceNodes); err != nil {
			global.Logger.Error("加载代理配置失败", zap.Error(err))
			global.ServiceManagerHandler.err = err
			return
		}

		lbManager := loadBalance.NewLoadBalance(vcfg.GetInt("load_type"))
		lbManager.Add(loadBalanceNodes...)

		global.ServiceManagerHandler.lb["http_proxy"] = lbManager
	})

	return global.ServiceManagerHandler.err
}
