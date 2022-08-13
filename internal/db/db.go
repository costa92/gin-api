package db

import (
	"github.com/rs/zerolog/log"
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
		log.Error().Msgf("init db err:%s", err)
		panic(err)
	}
}
