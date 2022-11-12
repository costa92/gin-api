package areas

import "gorm.io/gorm"

type AreaController struct {
	MysqlStorage *gorm.DB
}

func NewAreaController(db *gorm.DB) *AreaController {
	return &AreaController{
		MysqlStorage: db,
	}
}
