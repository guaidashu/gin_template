/**
  create by yy on 2019-07-29
*/

package data_struct

type Config struct {
	Mysql      MysqlConf
	PostGreSql PostGreSql
	App        AppConf
	Redis      RedisConf
}

type MysqlConf struct {
	Database   string `json:"database"`
	DbHost     string `json:"db_host"`
	DbPassword string `json:"db_password"`
	DbUsername string `json:"db_username"`
	DbPort     string `json:"db_port"`
	DbPoolSize int    `json:"db_pool_size"`
}

type AppConf struct {
	LogDir     string      `json:"log_dir"`
	RunAddress string      `json:"run_address"`
	RunPort    interface{} `json:"run_port"`
	DEBUG      bool        `json:"debug"`
}

type RedisConf struct {
	RedisHost     string `json:"redis_host"`
	RedisPort     string `json:"redis_port"`
	RedisPassword string `json:"redis_password"`
}

type PostGreSql struct {
	Database string `json:"database"`
	Host     string `json:"host"`
	Password string `json:"password"`
	Username string `json:"username"`
	Port     string `json:"port"`
	PoolSize int    `json:"pool_size"`
}
