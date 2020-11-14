package mysql

import (
	"Crd-End/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	mysqlConf = config.LoadMysql()
	url = mysqlConf["url"]
	username = mysqlConf["username"]
	password = mysqlConf["password"]
	Db, err = gorm.Open("mysql", username + ":" + password + "@" + url + "?charset=utf8&parseTime=True&loc=Local")
	//Db, err = gorm.Open("mysql", "root:123456@(202.107.190.8:10131)/blaststudio?charset=utf8&parseTime=True&loc=Local")
)

func InitMysql() {
	//defer Db.Close()
	Db.SingularTable(true)
	if err != nil {
		panic(err)
	}
}
