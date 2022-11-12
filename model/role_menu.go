package model

import (
	"context"

	"gorm.io/gorm"
)

const TableNameRoleMenu = "role_menu"

type RoleMenu struct {
	ID     uint `gorm:"column:id" json:"id"`
	RoleId uint `gorm:"column:role_id" json:"role_id"`
	MenuId uint `gorm:"column:menu_id" json:"menu_id"`
}

func (m *RoleMenu) TableName() string {
	return TableNameRoleMenu
}

type RoleMenuModel struct {
	DB *gorm.DB
}

func NewRoleMenuModel(ctx context.Context, db *gorm.DB) *RoleMenuModel {
	return &RoleMenuModel{
		DB: db.Model(RoleMenu{}).WithContext(ctx),
	}
}

func (m *RoleMenuModel) FindByRoleId(roleId int64) ([]*RoleMenu, error) {
	roleMenus := make([]*RoleMenu, 0)
	if err := m.DB.Where("role_id = ?", roleId).Find(&roleMenus).Error; err != nil {
		return nil, err
	}
	return roleMenus, nil
}

func (m *RoleMenuModel) FindByRoleIds(roleIds []int64) ([]*RoleMenu, error) {
	roleMenus := make([]*RoleMenu, 0)
	if err := m.DB.Where("role_id in ?", roleIds).Find(&roleMenus).Error; err != nil {
		return nil, err
	}
	return roleMenus, nil
}

func (m *RoleMenuModel) DeletedByRoleId(roleId int) error {
	return m.DB.Delete(&RoleMenu{}, "role_id = ?", roleId).Error
}

func (m *RoleMenuModel) GetRolesByRoleId(roleId int) ([]int, error) {
	var menuIds []int
	if err := m.DB.Select("menu_id").Where("role_id =?", roleId).Find(&menuIds).Error; err != nil {
		return nil, err
	}
	return menuIds, nil
}
