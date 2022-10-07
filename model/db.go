package model

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	lib "github.com/dappbujiujiu/sc_im/lib"
)

var Db *gorm.DB
func init() {
	dbConf := lib.GetConf().Db
	dsn := dbConf.User+":"+dbConf.Pass+"@tcp("+dbConf.Host+":"+dbConf.Port+")/"+dbConf.Name+"?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("db connect err:", err)
	}
	Db = db

	lib.GetConf()
}