package enterprise_type

import (
	"gorm.io/gorm"

	"github.com/costa92/go-web/model"
)

type EnterpriseTypeController struct {
	MysqlStorage *gorm.DB
}

func NewEnterpriseTypeController(db *gorm.DB) *EnterpriseTypeController {
	return &EnterpriseTypeController{MysqlStorage: db}
}

type CreateRequest struct {
	Name string `json:"name,omitempty"`
}

type UpdateRequest struct {
	ID int `json:"id" binding:"required"`
	CreateRequest
}

func (c *EnterpriseTypeController) SaveParams(enterpriseType *model.EnterpriseType, req *CreateRequest) {
	enterpriseType.Name = req.Name
}

type RequestDetail struct {
	ID int `query:"id" binding:"gte=1" form:"id"`
}
