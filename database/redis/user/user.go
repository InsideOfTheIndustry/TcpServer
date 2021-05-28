//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: user.go
// description: 与用户相关的redis实现
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-05-15
//

package user

import (
	"context"
	"strconv"
	"time"

	redisclient "github.com/InsideOfTheIndustry/TcpServe/database/redis"
	"github.com/InsideOfTheIndustry/TcpServe/logServer"

	"github.com/go-redis/redis/v8"
)

type UserCacheRepository struct {
	*redisclient.RedisEngine
}

var ctx = context.Background()

// SetVerificationCode(emailaddr, verificationcode string) error 实现存储库接口 设置验证码
func (urd UserCacheRepository) SetVerificationCode(emailaddr, verificationcode string) error {
	err := urd.Set(ctx, emailaddr, verificationcode, time.Duration(120*time.Second)).Err()
	if err != nil {
		logServer.Error("redis缓存设置验证码失败:%s", err.Error())
		return err
	}
	logServer.Info("设置验证码缓存成功。")
	return nil
}

// 	GetVerificationCode(emailaddr string)(string, error) 获取验证码
func (urd UserCacheRepository) GetVerificationCode(emailaddr string) (string, error) {
	VerificationCode, err := urd.Get(ctx, emailaddr).Result()
	if err == redis.Nil {
		logServer.Error("邮箱:(%s)的验证码不存在:(%s)", emailaddr, err.Error())
		return "", nil
	} else if err != nil {
		logServer.Error("读取发现错误:(%s)", err.Error())
		return "", err
	} else {
		logServer.Info("读取邮箱：(%s)的验证码成功", emailaddr)
		return VerificationCode, nil
	}
}

// SetToken 保存token数据
func (urd UserCacheRepository) SetToken(useraccount int64, token string) error {
	useraccountstring := strconv.FormatInt(useraccount, 10)
	err := urd.Set(ctx, useraccountstring, token, time.Duration(60*60*24*time.Second)).Err()
	if err != nil {
		logServer.Error("写入失败：%s", err.Error())
		return err
	}
	return nil
}

// GetToken(useraccount int64)(string, error)拿取token数据
func (urd UserCacheRepository) GetToken(useraccount int64) (string, error) {
	useraccountstring := strconv.FormatInt(useraccount, 10)
	tokenstring, err := urd.Get(ctx, useraccountstring).Result()
	if err == redis.Nil {
		logServer.Error("账户:(%s)的token不存在:(%s)", useraccountstring, err.Error())
		return "", nil
	} else if err != nil {
		logServer.Error("读取发现错误:(%s)", err.Error())
		return "", err
	} else {
		logServer.Info("读取账户:(%s)的token成功", useraccountstring)
		return tokenstring, nil
	}
}
