package distribute

import "gorm.io/gorm"

type DistributeController struct {
	MysqlStorage *gorm.DB
}

func NewDistributeController(db *gorm.DB) *DistributeController {
	return &DistributeController{MysqlStorage: db}
}
