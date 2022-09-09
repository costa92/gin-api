package auth

import "gorm.io/gorm"

type Auth struct {
	MysqlStorage *gorm.DB
}

func NewAuth(db *gorm.DB) *Auth {
	return &Auth{MysqlStorage: db}
}
