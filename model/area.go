package model

import (
	"context"

	"gorm.io/gorm"
)

var TableAreaName = "area"

type Area struct {
	ID    int64  `gorm:"column:id" json:"id"`
	Pid   int64  `gorm:"column:pid" json:"pid"`
	Name  string `gorm:"column:name" json:"name"`
	Level int    `gorm:"column:level" json:"level"`
}

// TableName Enterprise's table name
func (r *Area) TableName() string {
	return TableAreaName
}

type AreaModel struct {
	DB *gorm.DB
}

func NewAreaModel(ctx context.Context, tx *gorm.DB) *AreaModel {
	return &AreaModel{
		DB: tx.Model(&Area{}).WithContext(ctx),
	}
}

func (m *AreaModel) QueryProvinces() ([]*Area, error) {
	areas := make([]*Area, 0)
	if err := m.DB.Where("level = ?", 1).Debug().Find(&areas).Error; err != nil {
		return nil, err
	}
	return areas, nil
}

func (m *AreaModel) QueryListByPid(pid int64) ([]*Area, error) {
	areas := make([]*Area, 0)
	if err := m.DB.Where("pid = ?", pid).Find(&areas).Error; err != nil {
		return nil, err
	}
	return areas, nil
}

func (m *AreaModel) QueryListByIds(ids []int) ([]*Area, error) {
	areas := make([]*Area, 0)
	if err := m.DB.Where("id in (?)", ids).Order("id asc").Find(&areas).Debug().Error; err != nil {
		return nil, err
	}
	return areas, nil
}
