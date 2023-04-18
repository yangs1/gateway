package autoload

import (
	"fmt"
)

//MysqlConfig mysql信息配置
type MysqlConfig struct {
	Host         string `mapstructure:"host" json:"host" yaml:"host"`
	Port         int    `mapstructure:"port" json:"port" yaml:"port"`
	Name         string `mapstructure:"name" json:"name" yaml:"name"`
	Config       string `mapstructure:"config" json:"config" yaml:"config"`
	User         string `mapstructure:"user" json:"user" yaml:"user"`
	Password     string `mapstructure:"password" json:"password" yaml:"password"`
	Prefix       string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`
	MaxIdleConns int    `mapstructure:"max_idle_conns" json:"max_idle_conns" yaml:"max_open_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns" json:"max_open_conns" yaml:"max_open_conns"`
	MaxLifeTime  int    `mapstructure:"max_life_time" json:"max_life_time" yaml:"max_life_time"`
	LogMode      string `mapstructure:"log_mode" json:"log_mode" yaml:"log_mode"`
	LogZap       bool   `mapstructure:"log_zap" json:"log_zap" yaml:"log_zap"`
}

func (m *MysqlConfig) Dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", m.User, m.Password, m.Host, m.Port, m.Name, m.Config)
}
