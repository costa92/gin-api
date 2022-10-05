package db

import (
	"fmt"

	"github.com/costa92/errors"
	"gorm.io/gorm"

	"github.com/costa92/go-web/internal/option"
	"github.com/costa92/go-web/pkg/logger"
)

var (
	MysqlStorage *gorm.DB
	err          error
)

func GetMySQLFactoryOr(opts *option.MySQLOptions) (*gorm.DB, error) {
	if opts == nil && MysqlStorage == nil {
		return nil, fmt.Errorf("failed to get mysql store fatory")
	}
	MysqlStorage, err = NewClient(opts)
	if err != nil {
		logger.Errorf("init db err:%s", err)
		panic(err)
	}
	if MysqlStorage == nil || err != nil {
		return nil, fmt.Errorf("failed to get mysql store fatory, mysqlFactory: %+v, error: %w", MysqlStorage, err)
	}
	return MysqlStorage, nil
}

func Close() error {
	db, err := MysqlStorage.DB()
	if err != nil {
		return errors.Wrap(err, "get gorm db instance failed")
	}
	return db.Close()
}
