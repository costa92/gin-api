package model

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/costa92/go-web/pkg/meta"
)

const TableNameEnterpriseType = "enterprise_type"

// EnterpriseType mapped from table <enterprise_type>
type EnterpriseType struct {
	ID        int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"` // 主键
	Name      string `gorm:"column:name;not null" json:"name"`                  // 类型名称
	CreatedAt int64  `gorm:"column:created_at;not null" json:"created_at"`      // 添加时间
	UpdatedAt int64  `gorm:"column:updated_at;not null" json:"updated_at"`      // 修改时间
	UpdatedBy int32  `gorm:"column:updated_by;not null" json:"updated_by"`      // 修改者
	CreatedBy int32  `gorm:"column:created_by;not null" json:"created_by"`      // 添加者
}

// TableName EnterpriseType's table name
func (*EnterpriseType) TableName() string {
	return TableNameEnterpriseType
}

type EnterpriseTypeList struct {
	meta.ListMeta `json:",inline"`
	Items         []*EnterpriseType `json:"items"`
}

type EnterpriseTypeModel struct {
	DB *gorm.DB
}

func NewEnterpriseTypeModel(ctx context.Context, db *gorm.DB) *EnterpriseTypeModel {
	return &EnterpriseTypeModel{
		DB: db.Model(&EnterpriseType{}).WithContext(ctx),
	}
}

func (m *EnterpriseTypeModel) FirstByID(id int) (*EnterpriseType, error) {
	var enterpriseType EnterpriseType
	if err := m.DB.Where("id =?", id).First(&enterpriseType).Error; err != nil {
		return nil, err
	}
	return &enterpriseType, nil
}

func (m *EnterpriseTypeModel) Save(enterpriseType *EnterpriseType) error {
	tx := m.DB
	currTime := time.Now().Unix()
	enterpriseType.UpdatedAt = currTime
	if enterpriseType.ID > 0 {
		tx = tx.Where("id = ?", enterpriseType.ID)
	} else {
		enterpriseType.CreatedAt = currTime
	}

	if err := tx.Save(enterpriseType).Error; err != nil {
		return err
	}
	return nil
}

func (m *EnterpriseTypeModel) QueryList(name string) ([]*EnterpriseType, error) {
	tx := m.DB
	enterpriseTypes := make([]*EnterpriseType, 0)
	if name != "" {
		tx = tx.Where("name like ?", "%"+name+"%")
	}
	if err := tx.Find(&enterpriseTypes).Error; err != nil {
		return nil, err
	}
	return enterpriseTypes, nil
}
