/**
  create by yy on 2019-08-31
*/

package models

import (
	"fmt"
	"gin_template/app/config"
	"gin_template/app/data_struct/_interface"
	"gin_template/app/enum"
	"gin_template/app/libs"
	"gin_template/app/mongodb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

var GDB *gorm.DB

type (
	MysqlInit struct{}
)

func NewMysqlInit() *MysqlInit {
	return &MysqlInit{}
}

func (m *MysqlInit) Init(*_interface.ServiceParam) error {
	return InitDB()
}

func (m *MysqlInit) ComponentName() enum.BootModuleType {
	return enum.MysqlInit
}

func (m *MysqlInit) Close() error {
	if GDB != nil {
		libs.Logger.Info("Close Mysql")

		db, err := GDB.DB()
		if err != nil {
			return err
		}

		if err = db.Close(); err != nil {
			libs.Logger.Error("Close Mysql failed, error: %v", err)
			return err
		}
	}

	return nil
}

func InitDB() error {
	db, err := getDB()
	if err != nil {
		return err
	}
	GDB = db

	return nil
}

func getDB() (*gorm.DB, error) {
	connect := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local",
		config.Config.Mysql.DbUsername,
		config.Config.Mysql.DbPassword,
		config.Config.Mysql.DbHost,
		config.Config.Mysql.DbPort,
		config.Config.Mysql.Database)
	db, err := gorm.Open(mysql.Open(connect), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.Use(dbresolver.Register(dbresolver.Config{}).
		SetMaxIdleConns(config.Config.Mysql.DbPoolSize / 2).
		SetMaxOpenConns(config.Config.Mysql.DbPoolSize),
	)
	if err != nil {
		return nil, err
	}

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

	if mongodb.MDB != nil {
		libs.Logger.Info("Close Mongodb")
		if err = mongodb.MDB.Close(); err != nil {
			libs.Logger.Info("Close Mongodb failed, error: %v", err)
		}
	}
}
