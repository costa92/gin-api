package db

import (
	"gorm.io/gorm"

	"github.com/costa92/go-web/internal/option"
	"github.com/costa92/go-web/internal/pkg/db"
)

func NewClient(config *option.MySQLOptions) (*gorm.DB, error) {
	opt := &db.Options{
		Host:                  config.Addr,
		Username:              config.User,
		Password:              config.Pass,
		Database:              config.Database,
		MaxConnectionLifeTime: config.MaxConnectionLeftTime,
		MaxIdleConnections:    config.MaxIdleConnections,
		MaxOpenConnections:    config.MaxOpenConnections,
	}
	return db.New(opt)
}
