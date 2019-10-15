/**
  create by yy on 2019-08-31
*/

package models

import "github.com/jinzhu/gorm"

type UserModel struct {
	Model
	UserName string `gorm:"not null;size:255" json:"user_name"`
	Password string `gorm:"not null;size:255" json:"password"`
	Power    int    `gorm:"default:1" json:"power"`
}

func (u *UserModel) GetDB() *gorm.DB {
	db := GDB.Table(u.TableName())
	if db == nil {
		u.CreateTable()
		return GDB.Table(u.TableName()).Where("status = ?", 1)
	}
	return db
}

func (u *UserModel) TableName() string {
	return "user"
}

func (u *UserModel) CreateTable() {
	u.GetDB().Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(u)
}

func (u *UserModel) Create() {
	u.GetDB().Create(u)
}

func (u *UserModel) GetUserById(id int) (*[]UserModel, error) {
	var user []UserModel
	db := u.GetDB()
	err := db.Where("id=?", id).Find(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
