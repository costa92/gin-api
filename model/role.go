package model

import (
	"context"
	"gorm.io/gorm"

	"github.com/costa92/go-web/pkg/meta"
)

const TableNameRole = "roles"

type Role struct {
	ID        int64          `json:"id" gorm:"id"`
	Name      string         `gorm:"column:name" json:"name"`        // 角色名
	Remark    string         `gorm:"remark" json:"remark"`           // 备注
	Status    int            `gorm:"status;default:1" json:"status"` // 状态 1正常 2 禁止
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"-"`     // 删除时间
}

// TableName Role's table name
func (r *Role) TableName() string {
	return TableNameRole
}

type RoleList struct {
	meta.ListMeta `json:",inline"`
	Items         []*Role `json:"items"`
}

type RoleModel struct {
	DB *gorm.DB
}

func NewRoleModel(ctx context.Context, db *gorm.DB) *RoleModel {
	return &RoleModel{
		DB: db.Model(&Role{}).WithContext(ctx),
	}
}
