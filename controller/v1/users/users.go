package users

import "gorm.io/gorm"

type UserController struct {
	MysqlStorage *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{MysqlStorage: db}
}
