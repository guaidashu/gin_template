package models

import (
	"fmt"
	"gin_template/app/libs/serror"
	"sync"

	"gorm.io/gorm"
	"time"
)

var (
	templateCacheKey = "gin_template#template#cache#"
)

type (
	defaultTemplateModel struct {
		CacheConn
	}

	${TemplateStruct}
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
	return GDB.Migrator().HasTable(model.TableName())
}

func (model *defaultTemplateModel) CreateTable() error {
	db := model.getDB()
	err := db.Set("gorm:table_options", "ENGINE=InnoDB").Migrator().CreateTable(&TemplateModel{})
	if err != nil {
		return err
	}

	return nil
}

// 创建数据并返回本次插入的ID
func (model *defaultTemplateModel) Create(templateModel *TemplateModel) (err error) {
	db := model.getDB().Create(templateModel)
	if err = db.Error; err != nil {
		return
	}

	return
}

// 通过主键获取 数据
func (model *defaultTemplateModel) FindOne(Id int64) (*TemplateModel, error) {
	db := model.getDB().Where("id = ?", Id)

	templateModel := new(TemplateModel)
	key := fmt.Sprintf("%s%d", templateCacheKey, Id)
	err := model.QueryRow(templateModel, key, db, func(conn *gorm.DB, v interface{}) error {
		return conn.First(v).Error
	})

	return templateModel, err
}

// 单条更新, 多条更新请自行定义并维护键值
func (model *defaultTemplateModel) Update(templateModel *TemplateModel) error {
    key := fmt.Sprintf("%s%d", templateCacheKey, templateModel.Id)
    // 先删除一次缓存，防止redis操作失败导致数据不一致
	err := model.DelCache(key)
	if err != nil {
		return serror.NewErr().SetErr(err)
	}

	// 更新
	db := model.getDB()
	err := db.Where("id = ?", templateModel.Id).Save(templateModel).Error
	if err != nil {
		err = serror.NewErr().SetErr(err)
		return err
	}

	// 删除key
	return model.DelCache(key)
}

// 最小更新，每次只修改传入不为空或0的字段，如有需要改为空或0的，传入字段定义名
// 例：
// type modelName struct {
//     ExceptColumnsName1 string `json:"..."`
//     ExceptColumnsName1 int64  `json:"..."`
// }
// MinUpdate(templateModel, "ExceptColumnsName1", "ExceptColumnsName2")
func (model *defaultTemplateModel) MinUpdate(templateModel *TemplateModel, except ...string) error {
    key := fmt.Sprintf("%s%d", templateCacheKey, templateModel.Id)
    // 先删除一次缓存，防止redis操作失败导致数据不一致
	err := model.DelCache(key)
	if err != nil {
		return serror.NewErr().SetErr(err)
	}

	// 先转换为更新map
	update, err := struct2Map(templateModel, NewExcept(except...), nil)
	if err != nil {
		err = serror.NewErr().SetErr(err)
		return err
	}

	// 更新
	delete(update, "id")
	db := model.getDB()
	err = db.Where("id = ?", templateModel.Id).Updates(update).Error
	if err != nil {
		err = serror.NewErr().SetErr(err)
		return err
	}

	// 删除key
	return model.DelCache(key)
}

// 通过主键ID获取单条数据 (已经被软删除的数据)
func (model *defaultTemplateModel) GetTemplateIdById(Id int64) (templateModel *TemplateModel, err error) {
	templateModel, err = model.FindOne(Id)
	if err != nil {
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
