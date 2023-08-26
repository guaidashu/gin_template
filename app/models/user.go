package models

import (
	"fmt"
	"gin_template/app/libs/serror"
	"sync"

	"gorm.io/gorm"
)

var (
	userCacheKey = "user#cache#"
)

type (
	defaultUserModel struct {
		CacheConn
	}

	UserModel struct {
		// 主键
		UserId int64 `gorm:"primary_key;column:user_id;auto_increment" json:"user_id"`
		// 用户名
		Username string `gorm:"not null;size:255;column:username" json:"username"`
		// 手机号
		PhoneNumber string `gorm:"size:11;column:phone_number;comment:'商家手机号'" json:"phone_number"`
		// 小程序 open_id
		OpenId string `gorm:"size:255;column:open_id;comment:'小程序 open_id'" json:"open_id"`
		// 登录token，备用
		Token string `gorm:"type:text;column:token;comment:'登录token，备用'" json:"token"`
		// 用户头像
		AvatarUrl string `gorm:"type:text;column:avatar_url;comment:'用户头像'" json:"avatar_url"`
		// 国家
		Country string `gorm:"column:country;comment:'国家'" json:"country"`
		// 省份
		Province string `gorm:"column:province;comment:'省份'" json:"province"`
		// 城市
		City string `gorm:"column:city;comment:'城市'" json:"city"`
		// 语言
		Language string `gorm:"column:language;comment:'语言'" json:"language"`
		// 性别
		Sex uint64 `gorm:"column:sex;comment:'性别'" json:"sex"`
		// 邮箱
		Email string `gorm:"column:email;comment:'邮箱'" json:"email"`
		// 创建时间
		Created int64 `gorm:"column:created;comment:'创建时间'" json:"created"`
		// 更新时间
		Updated int64 `gorm:"column:updated;comment:'更新时间'" json:"updated"`
		// 删除时间
		Deleted int64 `gorm:"column:deleted;default:0;comment:'删除时间'" json:"deleted"`
	}
)

var (
	_userModel *defaultUserModel
	_userOnce  sync.Once
)

func NewUserModel() *defaultUserModel {
	_userOnce.Do(func() {
		_userModel = &defaultUserModel{
			NewCacheConn(),
		}
	})

	return _userModel
}

func (model *defaultUserModel) getDB() *gorm.DB {
	return getTable(GDB, model.TableName())
}

func (model *defaultUserModel) getDBWithNoDeleted() *gorm.DB {
	return getTableWithNoDeleted(GDB, model.TableName())
}

func (model *defaultUserModel) TableName() string {
	return "user"
}

func (model *defaultUserModel) HasTable() bool {
	return GDB.Migrator().HasTable(model.TableName())
}

func (model *defaultUserModel) CreateTable() error {
	db := model.getDB()
	err := db.Set("gorm:table_options", "ENGINE=InnoDB").Migrator().CreateTable(&UserModel{})
	if err != nil {
		return err
	}
	return nil
}

// 创建数据并返回本次插入的ID
func (model *defaultUserModel) Create(userModel *UserModel) (err error) {
	db := model.getDB().Create(userModel)
	if err = db.Error; err != nil {
		return
	}

	return
}

// 通过主键获取 数据
func (model *defaultUserModel) FindOne(userId int64) (*UserModel, error) {
	db := model.getDB().Where("user_id = ?", userId)

	userModel := new(UserModel)
	key := fmt.Sprintf("%s%d", userCacheKey, userId)
	err := model.QueryRow(userModel, key, db, func(conn *gorm.DB, v interface{}) error {
		return conn.First(v).Error
	})

	return userModel, err
}

// 单条更新, 多条更新请自行定义并维护键值
func (model *defaultUserModel) Update(userModel *UserModel) error {
	db := model.getDB()
	key := fmt.Sprintf("%s%d", userCacheKey, userModel.UserId)

	// 先转换为更新map
	update, err := struct2Map(userModel, nil, nil)
	if err != nil {
		err = serror.NewErr().SetErr(err)
		return err
	}

	// 更新
	delete(update, "user_id")
	err = db.Where("user_id = ?", userModel.UserId).Updates(update).Error
	if err != nil {
		err = serror.NewErr().SetErr(err)
		return err
	}

	// 删除key
	return model.DelCache(key)
}

// 通过主键ID获取单条数据 (已经被软删除的数据)
func (model *defaultUserModel) GetUserById(userId int64) (userModel *UserModel, err error) {
	userModel, err = model.FindOne(userId)
	if err != nil {
		return
	}

	if userModel.Deleted != 0 {
		err = gorm.ErrRecordNotFound
		return
	}

	return
}

// get user info by open id
func (model *defaultUserModel) GetUserByOpenId(openId string) (user *UserModel, err error) {
	db := model.getDB()

	user = new(UserModel)
	if err = db.Where("open_id = ?", openId).First(user).Error; err != nil {
		return
	}

	return
}

// 获取列表 (分页)
// func (model *defaultUserModel) GetUserList(req *requests.UserListReq) (
// 	list []*UserModel, err error) {
// 	db := model.getDB()
//
// 	if db, err = queryAutoWhere(db, req); err != nil {
// 		err = serror.NewErr().SetErr(err)
// 		return
// 	}
//
// 	err = db.Order(getOrderByStr(req.OrderField, "created", req.OrderType)).
// 		Offset(getOffset(req.Page, req.Size)).Limit(req.Size).
// 		Find(&list).Error
// 	if err != nil {
// 		err = serror.NewErr().SetErr(err)
// 		return
// 	}
//
// 	return
// }
