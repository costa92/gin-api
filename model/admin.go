package model

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/costa92/go-web/internal/pkg/meta"
)

const TableNameAdmin = "admin"

// Admin mapped from table <admin>
type Admin struct {
	meta.ObjectMeta `json:"metadata,omitempty"`
	Nickname        string         `gorm:"column:nickname" json:"nickname"`     // 用户昵称
	Password        string         `gorm:"column:password" json:"password"`     // 密码
	Name            string         `gorm:"column:name" json:"name"`             // 密码
	Salt            string         `gorm:"column:salt" json:"salt"`             // 加盐
	UpdatedAt       time.Time      `gorm:"column:updated_at" json:"updated_at"` // 更新时间
	CreatedAt       time.Time      `gorm:"column:created_at" json:"created_at"` // 保存时间
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"` // 删除时间
}

// TableName Admin's table name
func (a *Admin) TableName() string {
	return TableNameAdmin
}

type AdminList struct {
	meta.ListMeta `json:",inline"`
	Items         []*Admin `json:"items"`
}

type AdminModel struct {
	DB *gorm.DB
}

func NewAdminModel(ctx context.Context, db *gorm.DB) *AdminModel {
	return &AdminModel{
		DB: db.Model(&Admin{}).WithContext(ctx),
	}
}

func (a *AdminModel) FirstByName(name string) (*Admin, error) {
	admin := &Admin{}
	if err := a.DB.Where("name = ?", name).First(admin).Error; err != nil {
		return nil, err
	}
	return admin, nil
}
