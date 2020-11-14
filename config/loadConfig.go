package config

import (
	"github.com/Unknwon/goconfig"
)

func LoadMysql() map[string]string{
	config,err := goconfig.LoadConfigFile("config/conf/config.ini")
	if err!=nil {
		panic(err)
	}
	zzz , err := config.GetSection("mysql")
	return zzz
}

func LoadRedis() map[string]string{
	config,err := goconfig.LoadConfigFile("config/conf/config.ini")
	if err!=nil {
		panic(err)
	}
	zzz , err := config.GetSection("redis")
	return zzz
}