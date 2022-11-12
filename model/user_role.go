package model

import (
	"context"

	"gorm.io/gorm"
)

const TableNameUserRole = "user_role"

type UserRole struct {
	ID     int `gorm:"column:id" json:"id"`
	RoleId int `gorm:"column:role_id" json:"role_id"`
	UserID int `gorm:"column:user_id" json:"user_id"`
}

func (m *UserRole) TableName() string {
	return TableNameUserRole
}

type UserRoleModel struct {
	DB *gorm.DB
}

func NewUserRoleModel(ctx context.Context, db *gorm.DB) *UserRoleModel {
	return &UserRoleModel{
		DB: db.Model(UserRole{}).WithContext(ctx),
	}
}

func (m *UserRoleModel) FindByRoleId(roleId int64) ([]*UserRole, error) {
	roleMenus := make([]*UserRole, 0)
	if err := m.DB.Where("role_id = ?", roleId).Find(&roleMenus).Error; err != nil {
		return nil, err
	}
	return roleMenus, nil
}

func (m *UserRoleModel) FindByUserId(userId int64) ([]*UserRole, error) {
	roleMenus := make([]*UserRole, 0)
	if err := m.DB.Where("user_id = ?", userId).Find(&roleMenus).Error; err != nil {
		return nil, err
	}
	return roleMenus, nil
}

func (m *UserRoleModel) FindByUserIds(userId []int64) ([]*UserRole, error) {
	roleMenus := make([]*UserRole, 0)
	if err := m.DB.Where("user_id in ?", userId).Find(&roleMenus).Error; err != nil {
		return nil, err
	}
	return roleMenus, nil
}

func (m *UserRoleModel) DeletedByUserId(userId int) error {
	return m.DB.Delete(&UserRole{}, "user_id = ?", userId).Error
}

func (m *UserRoleModel) GetRolesByUserId(userId int) ([]int, error) {
	var roleIds []int
	if err := m.DB.Select("role_id").Where("user_id =?", userId).Find(&roleIds).Error; err != nil {
		return nil, err
	}
	return roleIds, nil
}
