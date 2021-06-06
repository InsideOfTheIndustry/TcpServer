package jwt

import (
	"fmt"
	"testing"

	"github.com/InsideOfTheIndustry/TcpServe/config"
	"github.com/InsideOfTheIndustry/TcpServe/logServer"
)

func TestParseToken(t *testing.T) {
	logServer.Setup("info")                    // 设置日志等级
	config.Setup("../../config/settings.yaml") // 读取配置文件
	InitSecretkey()
	fmt.Println(config.ApplicationConfig)
	fmt.Println(SECRETKEY)

	if SECRETKEY != "" {
		t.Fail()
	}
	tokenstring, err := GenarateToken(100009)
	if err != nil {
		t.Errorf("token生成失败：%s", err.Error())
	}

	_, err = ParseToken(tokenstring)
	if err != nil {
		t.Errorf("token解析失败:%s", err.Error())
	}

}
