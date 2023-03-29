package util

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
	"runtime"
)

var rootPath string

func GetRootPath() string {
	if rootPath == "" {
		_, filePath, _, _ := runtime.Caller(0)
		dirPath := path.Dir(filePath)

		rootPath = path.Dir(dirPath + "/../../")
	}

	return rootPath
}

// IsExist 判断文件或者目录是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func ReadConfigFile(fileName string) (*viper.Viper, error) {
	rootDir := GetRootPath()

	v := viper.New()
	v.SetConfigType("yaml")
	v.AddConfigPath(rootDir + "/config")
	v.SetConfigName(fileName)

	if err := v.ReadInConfig(); err != nil {
		log.Printf("读取 %s 配置文件失败", fileName)
		return nil, err
	}
	return v, nil
}
