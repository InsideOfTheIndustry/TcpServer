//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: config.go
// description: 配置文件读取
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-04-27
//

package tool

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// ParseConfig 解析日志文件
func ParseConfig(path string) (*ConfigStruct, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer jsonFile.Close()

	fileReader := bufio.NewReader(jsonFile)
	byteJsonFile, err := ioutil.ReadAll(fileReader)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(byteJsonFile, &_config)
	if err != nil {
		return nil, err
	}
	return _config, nil
}

// GetConfig 获取配置结构
func GetConfig() *ConfigStruct {
	return _config
}
