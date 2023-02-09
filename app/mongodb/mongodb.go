/**
  create by yy on 2019-10-22
*/

package mongodb

import (
	"errors"
	"fmt"
	"gin_template/app/config"
	"gin_template/app/data_struct/_interface"
	"gin_template/app/enum"
	"gin_template/app/libs"
	"github.com/guaidashu/go_mongodb_yy"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	MDBInit struct{}
)

func NewMDBInit() *MDBInit {
	return &MDBInit{}
}

func (m *MDBInit) Init(*_interface.ServiceParam) error {
	InitMongoDB()
	return nil
}

func (m *MDBInit) ComponentName() enum.BootModuleType {
	return enum.MongoInit
}

func (m *MDBInit) Close() error {
	if MDB != nil {
		libs.Logger.Info("Close MDB")
		if err := MDB.Close(); err != nil {
			libs.Logger.Info("Close MDB failed, error: %v", err)
			return nil
		}
	}

	return nil
}

var MDB *go_mongodb_yy.MDBPool

func InitMongoDB() {
	MDB = getConnect()
}

func getConnect() *go_mongodb_yy.MDBPool {
	// 设置连接池 数量
	if config.Config.Mongodb.PoolSize != 0 {
		go_mongodb_yy.MDBPoolSize = config.Config.Mongodb.PoolSize
	}

	applyUrl, err := getApplyUrl()
	if err != nil {
		libs.DebugPrint(libs.GetErrorString(libs.NewReportError(err)))
		return nil
	}

	return go_mongodb_yy.NewClient(go_mongodb_yy.ClientOpts{
		Uri: applyUrl,
		Opt: options.Client(),
	})
}

func getApplyUrl() (applyUrl string, err error) {
	if config.Config.Mongodb.Host == "" || config.Config.Mongodb.Port == "" {
		return "mongodb://localhost:27017/admin", libs.NewReportError(errors.New("mongodb error: nil host or nil port"))
	}
	if config.Config.Mongodb.Username == "" {
		applyUrl = fmt.Sprintf("mongodb://%v:%v/%v",
			config.Config.Mongodb.Host,
			config.Config.Mongodb.Port,
			config.Config.Mongodb.Database)
	} else {
		applyUrl = fmt.Sprintf("mongodb://%v:%v@%v:%v/%v",
			config.Config.Mongodb.Username,
			config.Config.Mongodb.Password,
			config.Config.Mongodb.Host,
			config.Config.Mongodb.Port,
			config.Config.Mongodb.Database)
	}
	return
}
