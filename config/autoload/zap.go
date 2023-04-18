package autoload

import (
	"go.uber.org/zap/zapcore"
	"strings"
)

type Zap struct {
	Level       string `mapstructure:"level" yaml:"level" json:"level"`                      // 日志打印级别 debug  info  warning  error
	Format      string `mapstructure:"format" yaml:"format" json:"format,omitempty"`         // 输出日志格式	logfmt, json
	LogPath     string `mapstructure:"log_path" yaml:"log_path" json:"log_path"`             // 输出日志文件路径
	FileName    string `mapstructure:"file_name" yaml:"file_name" json:"file_name"`          // 文件名称
	EncodeLevel string `mapstructure:"encode_level" yaml:"encode_level" json:"encode_level"` // 编码级
	MaxSize     int    `mapstructure:"max_size" yaml:"max_size" json:"max_size"`             // 【日志分割】单个日志文件最多存储量 单位(mb)
	MaxBackups  int    `mapstructure:"max_backups" yaml:"max_backups" json:"max_backups"`    // 【日志分割】日志备份文件最多数量
	MaxAge      int    `mapstructure:"max_age" yaml:"max_age" json:"max_age"`                // 日志保留时间，单位: 天 (day)
	ShowLine    bool   `mapstructure:"show_line" yaml:"show_line" json:"show_line"`          //显示行
	Compress    bool   `mapstructure:"compress" yaml:"compress" json:"compress"`             // 是否压缩日志
	SplitLevel  bool   `mapstructure:"split_level" yaml:"split_level" json:"split_level"`    // 是否拆分多级日志
	LogStdout   bool   `mapstructure:"log_stdout" yaml:"log_stdout" json:"log_stdout"`       // 是否输出到控制台
}

// TransportLevel 根据字符串转化为 zapcore.Level
func (z *Zap) TransportLevel() zapcore.Level {
	z.Level = strings.ToLower(z.Level)
	switch z.Level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.WarnLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.DebugLevel
	}
}

// ZapEncodeLevel 根据 EncodeLevel 返回 zapcore.LevelEncoder
func (z *Zap) ZapEncodeLevel() zapcore.LevelEncoder {
	switch {
	case z.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		return zapcore.LowercaseLevelEncoder
	case z.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		return zapcore.LowercaseColorLevelEncoder
	case z.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
		return zapcore.CapitalLevelEncoder
	case z.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		return zapcore.CapitalColorLevelEncoder
	default:
		return zapcore.LowercaseLevelEncoder
	}
}
