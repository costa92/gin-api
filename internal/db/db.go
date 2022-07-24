package db

import (
	"github.com/costa92/go-web/config"
	"gorm.io/gorm"
)

var (
	MysqlStorage *gorm.DB
	err          error
)

func InitDB(conf *config.Config) {
	MysqlStorage, err = NewClient(conf.MysqlConfig)
	if err != nil {
		panic(err)
	}
}
