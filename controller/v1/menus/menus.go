package menus

import (
	"gorm.io/gorm"
)

type MenuController struct {
	MysqlStorage *gorm.DB
}

func NewMenuController(db *gorm.DB) *MenuController {
	return &MenuController{
		MysqlStorage: db,
	}
}
