package follow

import "gorm.io/gorm"

type FollowController struct {
	MysqlStorage *gorm.DB
}

func NewFollowController(db *gorm.DB) *FollowController {
	return &FollowController{
		MysqlStorage: db,
	}
}
