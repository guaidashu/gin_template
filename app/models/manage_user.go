package models

import (
	"fmt"
	"gin_template/app/libs/serror"
	"sync"

	"gorm.io/gorm"
)

var (
	manageUserCacheKey = "manageUser#cache#"
)

type (
	defaultManageUserModel struct {
		CacheConn
	}

	ManageUserModel struct {
		// 主键
		ManageUserId int64 `gorm:"primary_key;column:manage_user_id;auto_increment" json:"manage_user_id"`
		// 用户名
		UserName string `gorm:"not null;size:255;column:username" json:"username"`
		// 密码
		Password string `gorm:"not null;size:255;column:password" json:"-"`
		// 小程序 open_id
		PhoneNumber string `gorm:"size:11;column:phone_number" json:"phone_number"`
		// 登录token，备用
		Token string `gorm:"type:text;column:token" json:"token"`
		// 头像
		AvatarUrl string `gorm:"type:text;column:avatar_url" json:"avatar_url"`
		// 性别
		Sex uint64 `gorm:"column:sex" json:"sex"`
		// 邮箱
		Email string `gorm:"column:email" json:"email"`
		// 创建时间
		Created int64 `gorm:"column:created;comment:'创建时间'" json:"created"`
		// 更新时间
		Updated int64 `gorm:"column:updated;comment:'更新时间'" json:"updated"`
		// 删除时间
		Deleted int64 `gorm:"column:deleted;default:0;comment:'删除时间'" json:"deleted"`
	}
)

var (
	_manageUserModel *defaultManageUserModel
	_manageUserOnce  sync.Once
)

func NewManageUserModel() *defaultManageUserModel {
	_manageUserOnce.Do(func() {
		_manageUserModel = &defaultManageUserModel{
			NewCacheConn(),
		}
	})

	return _manageUserModel
}

func (model *defaultManageUserModel) getDB() *gorm.DB {
	return getTable(GDB, model.TableName())
}

func (model *defaultManageUserModel) getDBWithNoDeleted() *gorm.DB {
	return getTableWithNoDeleted(GDB, model.TableName())
}

func (model *defaultManageUserModel) TableName() string {
	return "manage_user"
}

func (model *defaultManageUserModel) HasTable() bool {
	return GDB.Migrator().HasTable(model.TableName())
}

func (model *defaultManageUserModel) CreateTable() error {
	db := model.getDB()
	err := db.Set("gorm:table_options", "ENGINE=InnoDB").Migrator().CreateTable(&ManageUserModel{})
	if err != nil {
		return err
	}
	return nil
}

// 创建数据并返回本次插入的ID
func (model *defaultManageUserModel) Create(manageUserModel *ManageUserModel) (err error) {
	db := model.getDB().Create(manageUserModel)
	if err = db.Error; err != nil {
		return
	}

	return
}

// 通过主键获取 数据
func (model *defaultManageUserModel) FindOne(manageUserId int64) (*ManageUserModel, error) {
	db := model.getDB().Where("manage_user_id = ?", manageUserId)

	manageUserModel := new(ManageUserModel)
	key := fmt.Sprintf("%s%d", manageUserCacheKey, manageUserId)
	err := model.QueryRow(manageUserModel, key, db, func(conn *gorm.DB, v interface{}) error {
		return conn.First(v).Error
	})

	return manageUserModel, err
}

// 单条更新, 多条更新请自行定义并维护键值
func (model *defaultManageUserModel) Update(manageUserModel *ManageUserModel) error {
	db := model.getDB()
	key := fmt.Sprintf("%s%d", manageUserCacheKey, manageUserModel.ManageUserId)

	// 先转换为更新map
	update, err := struct2Map(manageUserModel, nil, nil)
	if err != nil {
		err = serror.NewErr().SetErr(err)
		return err
	}

	// 更新
	delete(update, "manage_user_id")
	err = db.Where("manage_user_id = ?", manageUserModel.ManageUserId).Updates(update).Error
	if err != nil {
		err = serror.NewErr().SetErr(err)
		return err
	}

	// 删除key
	return model.DelCache(key)
}

// 通过主键ID获取单条数据 (已经被软删除的数据)
func (model *defaultManageUserModel) GetManageUserById(manageUserId int64) (manageUserModel *ManageUserModel, err error) {
	manageUserModel, err = model.FindOne(manageUserId)
	if err != nil {
		return
	}

	if manageUserModel.Deleted != 0 {
		err = gorm.ErrRecordNotFound
		return
	}

	return
}

// 通过手机号获取管理员信息
func (model *defaultManageUserModel) GetByPhoneNumber(phoneNumber string) (manageUser *ManageUserModel, err error) {
	manageUser = new(ManageUserModel)
	err = model.getDB().Where("phone_number = ?", phoneNumber).First(manageUser).Error
	return
}

// 获取列表 (分页)
// func (model *defaultManageUserModel) GetManageUserList(req *requests.ManageUserListReq) (
// 	list []*ManageUserModel, err error) {
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
