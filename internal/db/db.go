package db

import (
	"gorm.io/gorm"

	"github.com/costa92/go-web/config"
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
