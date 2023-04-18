package autoload

type Http struct {
	Port           string   `mapstructure:"port" yaml:"port" json:"port"`                                     //监听地址
	ReadTimeout    int      `mapstructure:"read_timeout" yaml:"read_timeout" json:"read_timeout"`             // 读取超时时长
	WriteTimeout   int      `mapstructure:"write_timeout" yaml:"write_timeout" json:"write_timeout"`          // 写入超时时长
	MaxHeaderBytes int      `mapstructure:"max_header_bytes" yaml:"max_header_bytes" json:"max_header_bytes"` // 最大的header大小，二进制位长度
	TimeLoc        string   `mapstructure:"time_loc" yaml:"time_loc" json:"time_loc"`                         // 设置时区
	AllowIp        []string `mapstructure:"allow_ip" yaml:"allow_ip" json:"allow_ip"`                         // 白名单ip列表
}
