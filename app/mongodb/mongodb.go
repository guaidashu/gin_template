/**
  create by yy on 2019-10-22
*/

package mongodb

import (
	"errors"
	"fmt"
	"gin_template/app/config"
	"gin_template/app/libs"
)

func getApplyUrl() (applyUrl string, err error) {
	//"mongodb://localhost:27017"
	if config.Config.Mongodb.Host == "" || config.Config.Mongodb.Port == "" {
		return "mongodb://localhost:27017", libs.NewReportError(errors.New("mongodb error: nil host or nil port"))
	}
	if config.Config.Mongodb.Username == "" {
		applyUrl = fmt.Sprintf("mongodb://%v:%v", config.Config.Mongodb.Host, config.Config.Mongodb.Port)
	} else {
		applyUrl = fmt.Sprintf("mongodb://%v:%v@%v:%v", config.Config.Mongodb.Username, config.Config.Mongodb.Password, config.Config.Mongodb.Host, config.Config.Mongodb.Port)
	}
	return
}
