package models

import (
	"fmt"
	"gin_template/app/libs/serror"
	"sync"

	"github.com/jinzhu/gorm"
)

var (
	templateCacheKey = "template#cache#"
)

type (
	defaultTemplateModel struct {
		CacheConn
	}

	TemplateModel struct {
		// 主键
		TemplateIdId int64 `gorm:"primary_key;column:template_id_id;auto_increment;comment:'主键'" json:"template_id_id"`
		// 创建时间
        Created int64 `gorm:"column:created;comment:'创建时间'" json:"created"`
        // 更新时间
        Updated int64 `gorm:"column:updated;comment:'更新时间'" json:"updated"`
        // 删除时间
        Deleted int64 `gorm:"column:deleted;default:0;comment:'删除时间'" json:"deleted"`
	}
)

var (
	_templateModel *defaultTemplateModel
	_templateOnce  sync.Once
)

func NewTemplateModel() *defaultTemplateModel {
	_templateOnce.Do(func() {
		_templateModel = &defaultTemplateModel{
			NewCacheConn(),
		}
	})

	return _templateModel
}

func (model *defaultTemplateModel) getDB() *gorm.DB {
	return getTable(GDB, model.TableName())
}

func (model *defaultTemplateModel) getDBWithNoDeleted() *gorm.DB {
	return getTableWithNoDeleted(GDB, model.TableName())
}

func (model *defaultTemplateModel) TableName() string {
	return "template_name"
}

func (model *defaultTemplateModel) HasTable() bool {
	return GDB.HasTable(model.TableName())
}

func (model *defaultTemplateModel) CreateTable() error {
	db := model.getDB()
	err := db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&TemplateModel{}).Error
	if err != nil {
		return err
	}
	return nil
}

// 创建数据并返回本次插入的ID
func (model *defaultTemplateModel) Create(templateModel *TemplateModel) (templateIdId int64, err error) {
	db := model.getDB().Create(templateModel)
	if err = db.Error; err != nil {
		return
	}
	// 获取本次ID
	templateIdId = db.RowsAffected

	return
}

// 通过主键获取 数据
func (model *defaultTemplateModel) FindOne(templateIdId int64) (*TemplateModel, error) {
	db := model.getDB().Where("template_id_id = ?", templateIdId)

	templateModel := new(TemplateModel)
	key := fmt.Sprintf("%s%d", templateCacheKey, templateIdId)
	err := model.QueryRow(templateModel, key, db, func(conn *gorm.DB, v interface{}) error {
		return conn.First(v).Error
	})

	return templateModel, err
}

// 单条更新, 多条更新请自行定义并维护键值
func (model *defaultTemplateModel) Update(templateModel *TemplateModel) error {
	db := model.getDB()
	key := fmt.Sprintf("%s%d", templateCacheKey, templateModel.TemplateIdId)

	// 先转换为更新map
	update, err := struct2Map(templateModel, nil)
	if err != nil {
		err = serror.NewErr().SetErr(err)
		return err
	}

	// 更新
	delete(update, "template_id_id")
	err = db.Where("template_id_id = ?", templateModel.TemplateIdId).Updates(update).Error
	if err != nil {
		err = serror.NewErr().SetErr(err)
		return err
	}

	// 删除key
	return model.DelCache(key)
}

// 通过主键ID获取单条数据 (已经被软删除的数据)
func (model *defaultTemplateModel) GetTemplateIdById(templateIdId int64) (templateModel *TemplateModel, err error) {
	templateModel, err = model.FindOne(templateIdId)
	if err != nil {
		return
	}

	if templateModel.Deleted != 0 {
		err = gorm.ErrRecordNotFound
		return
	}

	return
}

// 获取列表 (分页)
// func (model *defaultTemplateModel) GetTemplateList(req *requests.TemplateListReq) (
// 	list []*TemplateModel, err error) {
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