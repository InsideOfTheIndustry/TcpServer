//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: jwt.go
// description: jwt验证
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-05-18
//

package jwt

import (
	"time"

	"github.com/InsideOfTheIndustry/TcpServe/config"
	"github.com/InsideOfTheIndustry/TcpServe/logServer"
	"github.com/dgrijalva/jwt-go"
)

var SECRETKEY = config.ApplicationConfig.JwtKey

// UserClaim 用户token格式
type UserClaim struct {
	UserAccount int64
	jwt.StandardClaims
}

// GenarateToken 生成token数据
func GenarateToken(useraccount int64) (string, error) {
	// 初始化一个自定义的claim
	var expire = 60 * 60 * 24
	userclaim := UserClaim{
		UserAccount: useraccount,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(expire) * time.Second).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userclaim)
	tokenstring, err := token.SignedString([]byte(SECRETKEY))
	if err != nil {
		logServer.Error("生成token失败：%s", err.Error())
		return "", err
	}
	logServer.Info("生成token成功。")
	return tokenstring, err

}

// ParseToken 解析token
func ParseToken(tokenstring string) (*UserClaim, error) {
	token, err := jwt.ParseWithClaims(tokenstring, &UserClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRETKEY), nil
	})
	if claims, ok := token.Claims.(*UserClaim); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
