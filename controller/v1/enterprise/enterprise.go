package enterprise

import (
	"gorm.io/gorm"

	"github.com/costa92/go-web/model"
)

type EnterpriseController struct {
	MysqlStorage *gorm.DB
}

func NewEnterpriseController(db *gorm.DB) *EnterpriseController {
	return &EnterpriseController{MysqlStorage: db}
}

type CreateRequest struct {
	Name       string `json:"name" from:"name"  validate:"required"`
	ProvinceId int    `json:"province_id" from:"province_id"  validate:"required"`
	CityId     int    `json:"city_id" from:"city_id"  validate:"required"`
	CountyId   int    `json:"county_id" from:"county_id"  validate:"required"`
	Status     int    `json:"status" from:"status"  validate:"required"`
	Type       int    `json:"type" from:"type"  validate:"required"`
	Tel        string `json:"tel" from:"tel"  validate:"required"`
	Fax        string `json:"fax" from:"fax"`
}

func (c *EnterpriseController) saveParams(enterprise *model.Enterprise, req *CreateRequest) {
	enterprise.Name = req.Name
	enterprise.ProvinceId = req.ProvinceId
	enterprise.CityId = req.CityId
	enterprise.CountyId = req.CountyId
	enterprise.Status = req.Status
	enterprise.Type = req.Type
	enterprise.Tel = req.Tel
	enterprise.Fax = req.Fax
}

type UpdateRequest struct {
	Id int `json:"id,omitempty" from:"id" validate:"required"`
	CreateRequest
}

type DetailRequest struct {
	Id int `query:"id" binding:"gte=1" form:"id"`
}

type DetailResponse struct {
	*model.Enterprise
	AreaId   []int                      `json:"parentCode"`
	AreaName []string                   `json:"area_name"`
	TypeName string                     `json:"type_name"`
	Contacts []*model.EnterpriseContact `json:"contacts"`
}

type UpdateStatusRequest struct {
	Id     int `json:"id,omitempty" query:"id" from:"id" validate:"required"`
	Status int `json:"status" query:"status" from:"id" validate:"required"`
}
