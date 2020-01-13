/**
  create by yy on 2019-08-31
*/

package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type BaseModel interface {
	TableName() string
	GetDB() *gorm.DB
	CreateTable() error
	HasTable() bool
}

type Model struct {
	Id        int        `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Status    int        `gorm:"default:1" json:"status"`
}

func getTable(gdb *gorm.DB, tableName string) *gorm.DB {
	db := gdb.Table(tableName)
	if db != nil {
		db = db.Where("status = ?", 1)
	}
	return db
}
