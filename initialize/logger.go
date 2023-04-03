package initialize

import (
	"gateway/global"
	"gateway/pkg/util"
	"github.com/mitchellh/mapstructure"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type LogConfigs struct {
	Level       string // 日志打印级别 debug  info  warning  error
	Format      string // 输出日志格式	logfmt, json
	LogPath     string // 输出日志文件路径
	LogFileName string // 输出日志文件名称
	MaxSize     int    // 【日志分割】单个日志文件最多存储量 单位(mb)
	MaxBackups  int    // 【日志分割】日志备份文件最多数量
	MaxAge      int    // 日志保留时间，单位: 天 (day)
	Compress    bool   // 是否压缩日志
	LogStdout   bool   // 是否输出到控制台
}

func InitLogger() error {
	vcfg, err := util.ReadConfigFile("app")
	if err != nil {
		return err
	}

	logcfg := LogConfigs{}

	if err := mapstructure.Decode(vcfg.Get("logConf"), &logcfg); err != nil {
		return err
	}

	writeSyncer, err := getLogWriter(logcfg) // 日志文件配置 文件位置和切割
	if err != nil {
		return err
	}
	encoder := getEncoder(logcfg) // 获取日志输出编码

	logLevel := map[string]zapcore.Level{
		"debug": zapcore.DebugLevel,
		"info":  zapcore.InfoLevel,
		"warn":  zapcore.WarnLevel,
		"error": zapcore.ErrorLevel,
	}

	level, ok := logLevel[logcfg.Level] // 日志打印级别
	if !ok {
		level = logLevel["info"]
	}

	core := zapcore.NewCore(encoder, writeSyncer, level)

	lg := zap.New(core, zap.AddCaller()) // zap.Addcaller() 输出日志打印文件和行数如： logger/logger_test.go:33
	// 1. zap.ReplaceGlobals 函数将当前初始化的 logger 替换到全局的 logger,
	// 2. 使用 logger 的时候 直接通过 zap.S().Debugf("xxx") or zap.L().Debug("xxx")
	// 3. 使用 zap.S() 和 zap.L() 提供全局锁，保证一个全局的安全访问logger的方式
	zap.ReplaceGlobals(lg)

	global.Logger = lg

	return nil
}

// getLogWriter 获取日志输出方式  日志文件 控制台
func getLogWriter(conf LogConfigs) (zapcore.WriteSyncer, error) {
	defaultLogPath := util.GetRootPath() + "/storage/logs/"
	workerLogPath := conf.LogPath
	dirLists := strings.Split(conf.LogPath, "/")

	if dirLists[0] == "." || dirLists[0] == ".." {
		workerLogPath = path.Dir(util.GetRootPath() + "/" + conf.LogPath + "/")
	}

	// 判断日志路径是否存在，如果不存在就创建
	if exist := util.IsExist(workerLogPath); !exist {
		if workerLogPath == "" {
			workerLogPath = defaultLogPath
		}

		if err := os.MkdirAll(workerLogPath, os.ModePerm); err != nil {
			workerLogPath = defaultLogPath
			if err := os.MkdirAll(workerLogPath, os.ModePerm); err != nil {
				return nil, err
			}
		}
	}

	// 日志文件 与 日志切割 配置
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filepath.Join(workerLogPath, conf.LogFileName), // 日志文件路径
		MaxSize:    conf.MaxSize,                                   // 单个日志文件最大多少 mb
		MaxBackups: conf.MaxBackups,                                // 日志备份数量
		MaxAge:     conf.MaxAge,                                    // 日志最长保留时间
		Compress:   conf.Compress,                                  // 是否压缩日志
	}

	if conf.LogStdout {
		// 日志同时输出到控制台和日志文件中
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(lumberJackLogger), zapcore.AddSync(os.Stdout)), nil
	} else {
		// 日志只输出到日志文件
		return zapcore.AddSync(lumberJackLogger), nil
	}
}

// getEncoder 编码器(如何写入日志)
func getEncoder(conf LogConfigs) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder   // log 时间格式 例如: 2021-09-11t20:05:54.852+0800
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // 输出level序列化为全大写字符串，如 INFO DEBUG ERROR
	//encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	//encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	if conf.Format == "json" {
		return zapcore.NewJSONEncoder(encoderConfig) // 以json格式写入
	}
	return zapcore.NewConsoleEncoder(encoderConfig) // 以logfmt格式写入
}
