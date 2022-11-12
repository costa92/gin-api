package model

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/costa92/go-web/pkg/meta"
)

const TableNameDistribute = "distribute"

// Distribute mapped from table <distribute>
type Distribute struct {
	ID           int32 `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`  // 主键
	EnterpriseID uint  `gorm:"column:enterprise_id;not null" json:"enterprise_id"` // 企业id
	UserID       int32 `gorm:"column:user_id;not null" json:"user_id"`             // 跟进者
	Status       int32 `gorm:"column:status;not null" json:"status"`               // 状态 1 正常  2 无效
	UpdatedAt    int64 `gorm:"column:updated_at;not null" json:"updated_at"`       // 修改时间
	UpdatedBy    int32 `gorm:"column:updated_by;not null" json:"updated_by"`       // 修改者
	CreatedAt    int64 `gorm:"column:created_at;not null" json:"created_at"`       // 添加时间
	CreatedBy    int32 `gorm:"column:created_by;not null" json:"created_by"`       // 添加者
}

// TableName Distributes table name
func (*Distribute) TableName() string {
	return TableNameDistribute
}

type DistributeModel struct {
	DB *gorm.DB
}

type DistributeList struct {
	meta.ListMeta `json:",inline"`
	Items         []*Distribute `json:"items"`
}

func NewDistributeModel(ctx context.Context, db *gorm.DB) *DistributeModel {
	return &DistributeModel{
		DB: db.Model(Distribute{}).WithContext(ctx),
	}
}

func (m *DistributeModel) FirstByEnterpriseID(enterpriseID uint) (*Distribute, error) {
	var distribute Distribute
	if err := m.DB.Where("enterprise_id =?", enterpriseID).
		Where("status = ?", 1).Debug().
		First(&distribute).Error; err != nil {
		return nil, err
	}
	return &distribute, nil
}

func (m *DistributeModel) Save(distribute *Distribute) error {
	tx := m.DB
	currTime := time.Now().Unix()
	distribute.UpdatedAt = currTime
	if distribute.ID > 0 {
		tx = tx.Where("id = ?", distribute.ID)
	} else {
		distribute.CreatedAt = currTime
	}
	if err := tx.Save(distribute).Error; err != nil {
		return err
	}
	return nil
}
