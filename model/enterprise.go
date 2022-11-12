package model

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/costa92/go-web/pkg/meta"
)

var TableEnterpriseName = "enterprise"

type Enterprise struct {
	Id         uint   `gorm:"column:id" json:"id"`
	Name       string `gorm:"column:name" json:"name"`
	ProvinceId int    `gorm:"column:province_id" json:"province_id"`
	CityId     int    `gorm:"column:city_id" json:"city_id"`
	CountyId   int    `gorm:"column:county_id" json:"county_id"`
	Status     int    `gorm:"column:status" json:"status"`
	Type       int    `gorm:"column:type" json:"type"`
	Tel        string `gorm:"column:tel" json:"tel"`
	Fax        string `gorm:"column:fax" json:"fax"`
	CreatedAt  int64  `gorm:"column:created_at" json:"created_at"`
	CreatedBy  int    `gorm:"column:created_by" json:"created_by"`
	UpdateAt   int64  `gorm:"column:updated_at" json:"updated_at"`
	UpdateBy   int    `gorm:"column:updated_by" json:"updated_by"`
}

// TableName Enterprise's table name
func (r *Enterprise) TableName() string {
	return TableEnterpriseName
}

type EnterpriseList struct {
	meta.ListMeta `json:",inline"`
	Items         []*Enterprise `json:"items"`
}
type EnterpriseModel struct {
	DB *gorm.DB
}

func NewEnterpriseModel(ctx context.Context, tx *gorm.DB) *EnterpriseModel {
	return &EnterpriseModel{
		DB: tx.Model(&Enterprise{}).WithContext(ctx),
	}
}

func (m *EnterpriseModel) FirstById(id int) (*Enterprise, error) {
	var enterprise Enterprise
	if err := m.DB.Where("id =?", id).First(&enterprise).Error; err != nil {
		return nil, err
	}
	return &enterprise, nil
}

func (m *EnterpriseModel) Save(enterprise *Enterprise) error {
	tx := m.DB
	currTime := time.Now().Unix()
	enterprise.UpdateAt = currTime
	if enterprise.Id > 0 {
		tx = tx.Where("id = ?", enterprise.Id)
	} else {
		enterprise.CreatedAt = currTime
	}

	if err := tx.Save(enterprise).Debug().Error; err != nil {
		return err
	}
	return nil
}

func (m *EnterpriseModel) FindByIds(id []uint) ([]*Enterprise, error) {
	enterprises := make([]*Enterprise, 0)
	if err := m.DB.Where("id in ?", id).Find(&enterprises).Debug().Error; err != nil {
		return nil, err
	}
	return enterprises, nil
}

func (m *EnterpriseModel) FindByArea(provinceId, cityId, countyId, status int) ([]*Enterprise, error) {
	enterprises := make([]*Enterprise, 0)
	tx := m.DB.Where("province_id = ?", provinceId).
		Where("city_id =?", cityId).
		Where("county_id = ?", countyId).
		Where("status =?", status)
	if err := tx.Find(&enterprises).Error; err != nil {
		return nil, err
	}
	return enterprises, nil
}

func (m *EnterpriseModel) CountByArea(provinceId, cityId, countyId, status int) (int64, error) {
	var total int64
	tx := m.DB.Where("province_id = ?", provinceId).
		Where("city_id =?", cityId).
		Where("county_id = ?", countyId).
		Where("status =?", status)
	if err := tx.Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}
