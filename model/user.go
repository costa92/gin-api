package model

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/costa92/go-web/pkg/logger"
	"github.com/costa92/go-web/pkg/meta"
	"github.com/costa92/go-web/pkg/util/auth"
)

const TableNameUser = "users"

// User mapped from table <admin>
type User struct {
	ID        int            `column:"id" json:"id" `
	Nickname  string         `gorm:"column:nickname" json:"nickname"`     // 用户昵称
	Password  string         `gorm:"column:password" json:"-"`            // 密码
	Username  string         `gorm:"column:username" json:"username"`     // 密码
	Salt      string         `gorm:"column:salt" json:"-"`                // 加盐
	LastTime  int64          `gorm:"column:last_time" json:"last_time"`   // 最后登录时间
	UpdatedAt int64          `gorm:"column:updated_at" json:"updated_at"` // 更新时间
	CreatedAt int64          `gorm:"column:created_at" json:"created_at"` // 保存时间
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"-"`          // 删除时间
}

// TableName User's table name
func (u *User) TableName() string {
	return TableNameUser
}

func (u *User) Compare(pwd string) error {
	if err := auth.Compare(u.Password, pwd); err != nil {
		logger.Errorf("auth.Compare failed", "err", err)
		return fmt.Errorf("failed to compile password: %w", err)
	}
	return nil
}

type UserList struct {
	meta.ListMeta `json:",inline"`
	Items         []*User `json:"items"`
}

type UserModel struct {
	DB *gorm.DB
}

func NewUserModel(ctx context.Context, db *gorm.DB) *UserModel {
	return &UserModel{
		DB: db.Model(&User{}).WithContext(ctx),
	}
}

func (a *UserModel) FirstByName(name string) (*User, error) {
	admin := &User{}
	if err := a.DB.Where("username = ?", name).First(admin).Error; err != nil {
		return nil, err
	}
	return admin, nil
}

func (a *UserModel) Save(user *User) error {
	tx := a.DB
	if user.ID > 0 {
		tx = tx.Where("id = ?", user.ID)
	}
	if err := tx.Save(user).Error; err != nil {
		return err
	}
	return nil
}
