package enterprise_contact

import (
	"encoding/json"

	"gorm.io/gorm"

	"github.com/costa92/go-web/model"
)

type EnterpriseContactController struct {
	MysqlStorage *gorm.DB
}

func NewEnterpriseContactController(db *gorm.DB) *EnterpriseContactController {
	return &EnterpriseContactController{MysqlStorage: db}
}

type CreateRequest struct {
	Name         string   `json:"name" form:"name"  validate:"required"`
	EnterpriseId uint     `json:"enterprise_id" form:"enterprise_id" validate:"required"`
	Mobile       string   `json:"mobile" form:"mobile" validate:"required"`
	Position     string   `json:"position" form:"position"`
	MultiMobile  []string `json:"[]multi_mobile" form:"[]multi_mobile"`
	Status       int      `json:"status"`
}

func (c *EnterpriseContactController) saveParams(contact *model.EnterpriseContact, req *CreateRequest) {
	contact.Name = req.Name
	contact.EnterpriseID = req.EnterpriseId
	contact.Mobile = req.Mobile
	contact.Position = req.Position
	contact.Status = req.Status
	var multiMobile string
	if len(req.MultiMobile) > 0 {
		multiMobileByte, _ := json.Marshal(req.MultiMobile)
		multiMobile = string(multiMobileByte)
	}
	contact.MultiMobile = multiMobile
}

func (c *EnterpriseContactController) updateParams(contact *model.EnterpriseContact, req *UpdateRequest) {
	contact.Name = req.Name
	contact.Mobile = req.Mobile
	contact.Position = req.Position
}

type UpdateRequest struct {
	Id       int    `json:"id" form:"id"  validate:"required"`
	Name     string `json:"name" form:"name"  validate:"required"`
	Mobile   string `json:"mobile" form:"mobile" validate:"required"`
	Position string `json:"position" form:"position"`
}
