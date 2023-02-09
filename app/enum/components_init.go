package enum

type BootModuleType string

const (
	NacosInit BootModuleType = "nacos"
	RedisInit BootModuleType = "redis"
	MysqlInit BootModuleType = "mysql"
	PsqlInit  BootModuleType = "postgresql"
	KafkaInit BootModuleType = "kafka"
	MongoInit BootModuleType = "mongo"
)
