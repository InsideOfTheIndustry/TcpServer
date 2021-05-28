//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: xorminit.go
// description: 初始化数据库引擎
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-05-15
//

package xormdatabase

import (
	"github.com/InsideOfTheIndustry/TcpServe/config"
	"github.com/InsideOfTheIndustry/TcpServe/logServer"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"xorm.io/core"
)

//XormEngine 数据库引擎结构体
type XormEngine struct {
	*xorm.Engine
}

// 全局数据库引擎
var DBEngine *XormEngine

// InitXormEngine 初始化数据库引擎
func InitXormEngine() error {

	connectexpression := config.DatabaseConfig.User + ":" + config.DatabaseConfig.Password + "@tcp(" + config.DatabaseConfig.Host + ":" + config.DatabaseConfig.Port + ")/" + config.DatabaseConfig.DBName + "?charset=" + config.DatabaseConfig.Charset // "root:888888@tcp(127.0.0.1:3306)/db_jpa_demo?charset=utf8"
	engine, err := xorm.NewEngine("mysql", connectexpression)
	if err != nil {
		logServer.Error("数据库连接失败：（%s）", err.Error())
		return err
	}
	engine.ShowSQL(config.DatabaseConfig.Showsql)

	engine.SetMapper(core.SameMapper{})
	// 分配空间并指向分配的空间
	var orm = new(XormEngine)
	orm.Engine = engine
	DBEngine = orm

	logServer.Info("mysql数据库连接成功")
	return nil
}
