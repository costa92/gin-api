package roles

import (
	"gorm.io/gorm"
)

type RoleController struct {
	MysqlStorage *gorm.DB
}

func NewRoleController(db *gorm.DB) *RoleController {
	return &RoleController{MysqlStorage: db}
}
