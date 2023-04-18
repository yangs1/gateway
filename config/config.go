package config

import (
	"gateway/config/autoload"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

//对应yaml文件结构
type ServiceConfig struct {
	HTTP  autoload.Http        `mapstructure:"http" json:"http" yaml:"http"`
	Zap   autoload.Zap         `mapstructure:"zap" json:"zap" yaml:"zap"`
	Mysql autoload.MysqlConfig `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
}

var (
	GVA_CONFIG ServiceConfig
	GVA_DB     *gorm.DB
	GVA_LOG    *zap.Logger
)
