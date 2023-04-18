package main

import (
	"fmt"
	"gateway/config"
	"gateway/initialize"
	"gateway/provider"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type UserBasic struct {
	Name          string
	PassWord      string
	Avatar        string
	Gender        string `gorm:"column:gender;default:male;type:varchar(6) comment 'male表示男,famale表示女'"` //gorm为数据库字段约束
	Phone         string `valid:"matches(^1[3-9]{1}\\d{9}$)"`                                            //valid为条件约束
	Email         string `valid:"email"`
	Identity      string
	ClientIp      string `valid:"ipv4"`
	ClientPort    string
	Salt          string     //盐值
	LoginTime     *time.Time `gorm:"column:login_time"`
	HeartBeatTime *time.Time `gorm:"column:heart_beat_time"`
	LoginOutTime  *time.Time `gorm:"column:login_out_time"`
	IsLoginOut    bool
	DeviceInfo    string //登录设备
}

func main() {
	initialize.Viper(&config.GVA_CONFIG)

	fmt.Println(config.GVA_CONFIG)

	config.GVA_LOG = initialize.Zap()
	config.GVA_DB = initialize.GormMysql()
	//config.GVA_LOG.Info("test", zap.String("sad", "dasd"))
	//config.GVA_LOG.Error("test", zap.String("sad", "dasd"))
	//fmt.Println(config.GVA_DB.Debug().Table("relations").Where("id= ?", 1).Select("*"))

	config.GVA_DB.Table("relations").Where("id = ? and owner_id=?", "23", "1").Rows()
	return

	initialize.InitLoadBalance()

	engine := provider.NewGatewayEngine()

	go func() {
		engine.StartHttpServer()
	}()
	// 当前的goroutine等待信号量
	quit := make(chan os.Signal)
	// 监控信号：SIGINT, SIGTERM, SIGQUIT
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// 这里会阻塞当前goroutine等待信号
	<-quit

	engine.StopHttpServer()

}
