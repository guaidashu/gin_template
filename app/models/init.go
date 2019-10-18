/**
  create by yy on 2019-08-31
*/

package models

import (
	"fmt"
	"gin_template/app/config"
	"gin_template/app/libs"
	"gin_template/app/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var GDB *gorm.DB
var PDB *gorm.DB

func InitDB() error {
	db, err := getDB()
	if err != nil {
		return err
	}
	GDB = db
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "doratv_" + defaultTableName
	}
	return nil
}

func InitPostGreDB() error {
	db, err := getPostGreDB()
	if err != nil {
		return err
	}
	PDB = db
	return nil
}

func getPostGreDB() (*gorm.DB, error) {
	connect := fmt.Sprintf("host=%v user=%v dbname=%v sslmode=disable password=%v",
		config.Config.PostGreSql.Host,
		config.Config.PostGreSql.Username,
		config.Config.PostGreSql.Database,
		config.Config.PostGreSql.Password,
	)
	db, err := gorm.Open("postgres", connect)
	if err != nil {
		return nil, err
	}
	db.DB().SetMaxIdleConns(config.Config.PostGreSql.PoolSize / 2)
	db.DB().SetMaxOpenConns(config.Config.PostGreSql.PoolSize)
	return db, nil
}

func getDB() (*gorm.DB, error) {
	connect := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local",
		config.Config.Mysql.DbUsername,
		config.Config.Mysql.DbPassword,
		config.Config.Mysql.DbHost,
		config.Config.Mysql.DbPort,
		config.Config.Mysql.Database)
	db, err := gorm.Open("mysql", connect)
	if err != nil {
		return nil, err
	}
	db.DB().SetMaxIdleConns(config.Config.Mysql.DbPoolSize / 2)
	db.DB().SetMaxOpenConns(config.Config.Mysql.DbPoolSize)
	return db, nil
}

func CreateTable() {
	user := new(UserModel)
	user.CreateTable()
}

func CloseDB() {
	var err error
	libs.Logger.Info("Close Postgresql.")
	if PDB != nil {
		if err = PDB.Close(); err != nil {
			libs.Logger.Error("Close Postgresql failed, error: %v", err)
		}
	}

	libs.Logger.Info("Close Mysql")
	if GDB != nil {
		if err = GDB.Close(); err != nil {
			libs.Logger.Error("Close Mysql failed, error: %v", err)
		}
	}

	libs.Logger.Info("Close redis")
	if redis.Redis != nil {
		if err = redis.Redis.Close(); err != nil {
			libs.Logger.Info("Close redis failed, error: %v", err)
		}
	}
}
