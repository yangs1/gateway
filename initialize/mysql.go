package initialize

import (
	"fmt"
	"gateway/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

type writer struct {
	logger.Writer
}

// Printf 格式化打印日志
func (w *writer) Printf(message string, data ...interface{}) {
	//fmt.Println("???????")
	if config.GVA_CONFIG.Mysql.LogZap {
		config.GVA_LOG.Info(fmt.Sprintf(message+"\n", data...))
	} else {
		w.Writer.Printf(message, data...)
	}
}

// 获取mysql 配置
func GormMysql() *gorm.DB {
	m := config.GVA_CONFIG.Mysql
	if m.Host == "" {
		return nil
	}

	mysqlConfig := mysql.Config{
		DSN:                       m.Dsn(), // DSN data source name
		DefaultStringSize:         255,     // string 类型字段的默认长度
		SkipInitializeWithVersion: false,   // 根据版本自动配置
		DisableDatetimePrecision:  true,    // 禁止 datatime 精度
		DontSupportRenameColumn:   true,
		DontSupportRenameIndex:    true,
	}

	var _gorm = new(_gorm)

	if db, err := gorm.Open(mysql.New(mysqlConfig), _gorm.getOrmConfig(m.Prefix, m.LogMode)); err != nil {
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		sqlDB.SetConnMaxLifetime(time.Duration(m.MaxLifeTime) * time.Second)
		return db
	}
}

type _gorm struct{}

func (g *_gorm) getOrmConfig(prefix string, logModel string) *gorm.Config {

	cfg := &gorm.Config{
		// 命名策略
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   prefix, // 表前缀，在表名前添加前缀，如添加用户模块的表前缀 user_
			SingularTable: false,  // 是否使用单数形式的表名，如果设置为 true，那么 User 模型会对应 users 表
		},

		DisableForeignKeyConstraintWhenMigrating: true,
	}

	var logWriter *writer // 自定义 log writer

	logWriter = &writer{log.New(os.Stdout, "\r\n", log.LstdFlags)}

	_default := logger.New(logWriter, logger.Config{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      logger.Warn,
		Colorful:      true,
	})

	if logModel == "debug" {
		cfg.Logger = _default.LogMode(logger.Info)
	}
	//cfg.Logger = logger.Default.LogMode(logger.Info) //_default.LogMode(logger.Info)

	return cfg
}
