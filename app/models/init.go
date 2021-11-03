/**
  create by yy on 2019-08-31
*/

package models

import (
	"fmt"
	"gin_template/app/config"
	"gin_template/app/libs"
	"gin_template/app/mongodb"
	"gin_template/app/rds"
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
		return "table_" + defaultTableName
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

func createTable(models ...BaseModel) {

	var (
		err error
	)

	for _, model := range models {
		if !model.HasTable() {
			if err = model.CreateTable(); err != nil {
				libs.Logger.Error(fmt.Sprintf("create table error: %v", libs.NewReportError(err)))
			}
		}
	}

}

func CreateTable() {
	createTable(
		NewUserModel(),
	)
}

func CloseDB() {
	var err error
	if PDB != nil {
		libs.Logger.Info("Close Postgresql.")
		if err = PDB.Close(); err != nil {
			libs.Logger.Error("Close Postgresql failed, error: %v", err)
		}
	}

	if GDB != nil {
		libs.Logger.Info("Close Mysql")
		if err = GDB.Close(); err != nil {
			libs.Logger.Error("Close Mysql failed, error: %v", err)
		}
	}

	if rds.Redis != nil {
		libs.Logger.Info("Close rds")
		if err = rds.Redis.Close(); err != nil {
			libs.Logger.Info("Close rds failed, error: %v", err)
		}
	}

	if mongodb.MDB != nil {
		libs.Logger.Info("Close Mongodb")
		if err = mongodb.MDB.Close(); err != nil {
			libs.Logger.Info("Close Mongodb failed, error: %v", err)
		}
	}
}
