package utils

import (
	"os"
)

//@function: DirIsExists
//@description: 文件目录是否存在
//@param: path string
//@return: bool, error

func DirIsExists(path string, withMk bool) (bool, error) {
	fi, err := os.Stat(path)
	if err == nil && fi.IsDir() {
		return true, nil
	}

	return false, err
}
