package model

import (
	"context"
	"time"

	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	code2 "github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/meta"
)

const TableNameMenu = "menus"

type Menu struct {
	ID            uint   `gorm:"column:id" json:"id"`
	ParentID      uint   `gorm:"column:parent_id" json:"parent_id"`
	MenuName      string `gorm:"column:menu_name" json:"menu_name"`
	MenuType      uint   `gorm:"column:menu_type" json:"menu_type"`
	MenuSort      uint   `gorm:"column:menu_sort" json:"menu_sort"`
	MenuStatus    int    `gorm:"column:menu_status" json:"menu_status"`
	HideMenu      int    `gorm:"column:hide_menu" json:"hide_menu"`
	ExternalLink  int    `gorm:"column:external_link" json:"external_link"`
	Permission    string `gorm:"column:permission" json:"permission"`
	ComponentName string `gorm:"column:component_name" json:"component_name"`
	Component     string `gorm:"column:component" json:"component"`
	KeepAlive     uint   `gorm:"column:keepalive" json:"keep_alive"`
	Icon          string `gorm:"column:icon" json:"icon"`
	Path          string `gorm:"column:path" json:"path"`
	CreatedBy     uint   `bson:"created_by" gorm:"column:created_by" json:"created_by"`
	CreatedAt     int64  `bson:"created_at" gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedBy     uint   `bson:"updated_by" gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     int64  `bson:"updated_at" gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (m *Menu) TableName() string {
	return TableNameMenu
}

type MenuList struct {
	meta.ListMeta `json:",inline"`
	Items         []*Menu `json:"items"`
}

type MenuModel struct {
	DB *gorm.DB
}

func NewMenuModel(ctx context.Context, db *gorm.DB) *MenuModel {
	return &MenuModel{
		DB: db.Model(&Menu{}).WithContext(ctx),
	}
}

func (m *MenuModel) FirstByID(id int) (*Menu, error) {
	menu := Menu{}
	if err := m.DB.Where("id = ?", id).First(&menu).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code2.ErrMenuNotFound, err.Error())
		}
		return nil, errors.WithCode(code2.ErrDatabase, err.Error())
	}
	return &menu, nil
}

func (m *MenuModel) Save(ctx *gin.Context, menu *Menu) error {
	tx := m.DB
	// operater := ctx.GetString(middleware.UsernameKey)
	currTime := time.Now().Unix()
	if menu.ID > 0 {
		tx = tx.Where("id =?", menu.ID)
		menu.UpdatedAt = currTime
	}
	menu.CreatedAt = currTime
	if err := tx.Save(menu).Error; err != nil {
		return errors.WithCode(code2.ErrDatabase, err.Error())
	}
	return nil
}
