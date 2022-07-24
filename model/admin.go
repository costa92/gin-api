package model

import (
	"time"

	"github.com/costa92/go-web/internal/pkg/meta"
	"gorm.io/gorm"
)

const TableNameAdmin = "admin"

// Admin mapped from table <admin>
type Admin struct {
	meta.ObjectMeta `json:"metadata,omitempty"`
	Nickname        string         `gorm:"column:nickname" json:"nickname"`     // 用户昵称
	Password        string         `gorm:"column:password" json:"password"`     // 密码
	Salt            string         `gorm:"column:salt" json:"salt"`             // 加盐
	UpdatedAt       time.Time      `gorm:"column:updated_at" json:"updated_at"` // 更新时间
	CreatedAt       time.Time      `gorm:"column:created_at" json:"created_at"` // 保存时间
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"` // 删除时间
}

// TableName Admin's table name
func (*Admin) TableName() string {
	return TableNameAdmin
}

type AdminList struct {
	meta.ListMeta `json:",inline"`
	Items         []*Admin `json:"items"`
}
