package service

import (
	"devplat/src/log"
	"os"
)

func init() {
	InitChaincodeDir()
}

func TouchFile(filename, content string) {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		// 创建文件失败处理
		log.Logger.Error(err.Error())
	} else {
		_, err = f.Write([]byte(content))
		if err != nil {
			// 写入失败处理
			log.Logger.Error(err.Error())
		}
	}
}

func InitChaincodeDir() {
	os.Mkdir("chaincode", os.ModePerm)
}
