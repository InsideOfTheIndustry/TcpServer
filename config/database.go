package config

import (
	"github.com/spf13/viper"
)

//Database 数据库配置文件
type Database struct {
	Type     string
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
	Charset  string
	Showsql  bool
}

func InitDatabase(cfg *viper.Viper) *Database {
	return &Database{
		Type:     cfg.GetString("type"),
		User:     cfg.GetString("user"),
		Password: cfg.GetString("password"),
		Host:     cfg.GetString("host"),
		Port:     cfg.GetString("port"),
		DBName:   cfg.GetString("dbname"),
		Charset:  cfg.GetString("charset"),
		Showsql:  cfg.GetBool("showsql"),
	}
}

var DatabaseConfig = new(Database)

// Redis缓存数据库
type Redis struct {
	Addr     string
	Password string
	Db       int
}

func InitRedis(cfg *viper.Viper) *Redis {
	return &Redis{
		Addr:     viper.GetString("addr"),
		Password: viper.GetString("password"),
		Db:       viper.GetInt("db"),
	}
}

var RedisConfig = new(Redis)
