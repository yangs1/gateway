package initialize

import (
	"fmt"
	"gateway/config"
	"gateway/utils"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"time"
)

func Zap() (logger *zap.Logger) {
	workDir := path.Dir(config.GVA_CONFIG.Zap.LogPath)
	if ok, _ := utils.DirIsExists(workDir, true); !ok { // 判断是否有Director文件夹
		fmt.Printf("create %v directory\n", workDir)
		_ = os.Mkdir(workDir, os.ModePerm)
	}

	var z = new(_zap)
	cores := z.getZapCores()
	logger = zap.New(zapcore.NewTee(cores...))

	if config.GVA_CONFIG.Zap.ShowLine { // zap.Addcaller() 输出日志打印文件和行数如： logger/logger_test.go:33
		logger = logger.WithOptions(zap.AddCaller())
	}

	// 1. zap.ReplaceGlobals 函数将当前初始化的 logger 替换到全局的 logger,
	// 2. 使用 logger 的时候 直接通过 zap.S().Debugf("xxx") or zap.L().Debug("xxx")
	// 3. 使用 zap.S() 和 zap.L() 提供全局锁，保证一个全局的安全访问logger的方式
	zap.ReplaceGlobals(logger)

	return logger
}

// getLogWriter 获取日志输出方式  日志文件 控制台
func getLogWriter(level string) (zapcore.WriteSyncer, error) {

	zapConfig := config.GVA_CONFIG.Zap

	fileName := fmt.Sprintf("%s_%s.log", config.GVA_CONFIG.Zap.FileName, level)
	if level == "" {
		fileName = fmt.Sprintf("%s.log", config.GVA_CONFIG.Zap.FileName)
	}

	// 日志文件 与 日志切割 配置
	lumberJackLogger := &lumberjack.Logger{
		Filename:   path.Join(zapConfig.LogPath, fileName), // 日志文件路径
		MaxSize:    zapConfig.MaxSize,                      // 单个日志文件最大多少 mb
		MaxBackups: zapConfig.MaxBackups,                   // 日志备份数量
		MaxAge:     zapConfig.MaxAge,                       // 日志最长保留时间
		Compress:   zapConfig.Compress,                     // 是否压缩日志
	}

	if zapConfig.LogStdout {
		// 日志同时输出到控制台和日志文件中
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(lumberJackLogger), zapcore.AddSync(os.Stdout)), nil
	} else {
		// 日志只输出到日志文件
		return zapcore.AddSync(lumberJackLogger), nil
	}
}

type _zap struct{}

// getEncoder 编码器(如何写入日志)
func (z *_zap) getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = z.CustomTimeEncoder                     // log 时间格式 例如: 2021-09-11t20:05:54.852+0800
	encoderConfig.EncodeLevel = config.GVA_CONFIG.Zap.ZapEncodeLevel() // 输出level序列化为全大写字符串，如 INFO DEBUG ERROR
	//encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	//encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	if config.GVA_CONFIG.Zap.Format == "json" {
		return zapcore.NewJSONEncoder(encoderConfig) // 以json格式写入
	}
	return zapcore.NewConsoleEncoder(encoderConfig) // 以logfmt格式写入
}

// CustomTimeEncoder 自定义日志输出时间格式
func (z *_zap) CustomTimeEncoder(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	// zapcore.ISO8601TimeEncoder  // log 时间格式 例如: 2021-09-11t20:05:54.852+0800
	encoder.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// getZapCores 根据配置文件的Level获取 []zapcore.Core
func (z *_zap) getZapCores() []zapcore.Core {
	cores := make([]zapcore.Core, 0)

	// 判断是否根据 level 拆分成多个 文件
	if config.GVA_CONFIG.Zap.SplitLevel {
		for level := config.GVA_CONFIG.Zap.TransportLevel(); level <= zapcore.FatalLevel; level++ {
			cores = append(cores, z.getEncoderCoreByLevel(level))
		}
	} else {
		cores = append(cores, z.getEncoderCore())
	}

	return cores
}

// getEncoderCoreByLevel 获取Encoder的 zapcore.Core
func (z *_zap) getEncoderCoreByLevel(level zapcore.Level) zapcore.Core {
	writer, err := getLogWriter(level.String()) // 使用lumberJackLogger进行日志分割
	if err != nil {
		fmt.Printf("Get Write Syncer Failed err:%v", err.Error())
		return nil
	}

	return zapcore.NewCore(z.getEncoder(), writer, z.getLevelPriority(level))
}

// getEncoderCore 获取Encoder的 zapcore.Core
func (z *_zap) getEncoderCore() zapcore.Core {
	writer, err := getLogWriter("") // 使用lumberJackLogger进行日志分割
	if err != nil {
		fmt.Printf("Get Write Syncer Failed err:%v", err.Error())
		return nil
	}

	level, err := zapcore.ParseLevel(config.GVA_CONFIG.Zap.Level)
	if err != nil {
		level = zap.InfoLevel
	}
	return zapcore.NewCore(z.getEncoder(), writer, level)

}

// 自定义错误级别范围
func (z *_zap) getLevelPriority(l zapcore.Level) zap.LevelEnablerFunc {
	switch l {
	case zapcore.DebugLevel:
		return func(level zapcore.Level) bool { // 调试级别
			return level == zap.DebugLevel
		}
	case zapcore.InfoLevel:
		return func(level zapcore.Level) bool { // 日志级别
			return level == zap.InfoLevel
		}
	case zapcore.WarnLevel:
		return func(level zapcore.Level) bool { // 警告级别
			return level == zap.WarnLevel
		}
	case zapcore.ErrorLevel:
		return func(level zapcore.Level) bool { // 错误级别
			return level == zap.ErrorLevel
		}
	case zapcore.DPanicLevel:
		return func(level zapcore.Level) bool { // dpanic级别
			return level == zap.DPanicLevel
		}
	case zapcore.PanicLevel:
		return func(level zapcore.Level) bool { // panic级别
			return level == zap.PanicLevel
		}
	case zapcore.FatalLevel:
		return func(level zapcore.Level) bool { // 终止级别
			return level == zap.FatalLevel
		}
	default:
		return func(level zapcore.Level) bool { // 调试级别
			return level == zap.DebugLevel
		}
	}
}
