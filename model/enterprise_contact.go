package model

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/costa92/go-web/pkg/meta"
)

const TableNameEnterpriseContact = "enterprise_contact"

const (
	_ = iota
	ContactStatusNormal
	ContactStatusFail
)

// EnterpriseContact mapped from table <enterprise_contact>
type EnterpriseContact struct {
	ID           int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`  // 主键
	Name         string `gorm:"column:name;not null" json:"name"`                   // 姓名
	EnterpriseID uint   `gorm:"column:enterprise_id;not null" json:"enterprise_id"` // 企业id
	MultiMobile  string `gorm:"column:multi_mobile;not null" json:"multi_mobile"`   // 多个手机号码
	Mobile       string `gorm:"column:mobile;not null" json:"mobile"`               // 联系电话
	Position     string `gorm:"column:position;not null" json:"position"`           // 职位
	Status       int    `gorm:"column:status;" json:"status"`                       // 状态 1 正常 2 失效
	UpdatedAt    int64  `gorm:"column:updated_at;not null" json:"updated_at"`       // 修改时间
	UpdatedBy    int    `gorm:"column:updated_by;not null" json:"updated_by"`       // 修改者
	CreatedAt    int64  `gorm:"column:created_at;not null" json:"created_at"`       // 添加时间
	CreatedBy    int    `gorm:"column:created_by;not null" json:"created_by"`       // 添加者
}

// TableName EnterpriseContact's table name
func (*EnterpriseContact) TableName() string {
	return TableNameEnterpriseContact
}

type EnterpriseContactList struct {
	meta.ListMeta `json:",inline"`
	Items         []*EnterpriseContact `json:"items"`
}

type EnterpriseContactModel struct {
	DB *gorm.DB
}

func NewEnterpriseContactModel(ctx context.Context, tx *gorm.DB) *EnterpriseContactModel {
	return &EnterpriseContactModel{
		DB: tx.Model(&EnterpriseContact{}).WithContext(ctx),
	}
}

func (m *EnterpriseContactModel) FirstById(id int) (*EnterpriseContact, error) {
	var contact EnterpriseContact
	if err := m.DB.Where("id =?", id).First(&contact).Error; err != nil {
		return nil, err
	}
	return &contact, nil
}

func (m *EnterpriseContactModel) FindByEnterpriseId(enterpriseId int) ([]*EnterpriseContact, error) {
	contacts := make([]*EnterpriseContact, 0)
	if err := m.DB.Where("enterprise_id = ?", enterpriseId).
		Where("status =?", ContactStatusNormal).Find(&contacts).Error; err != nil {
		return nil, err
	}
	return contacts, nil
}

func (m *EnterpriseContactModel) Save(contact *EnterpriseContact) error {
	tx := m.DB
	currTime := time.Now().Unix()
	contact.UpdatedAt = currTime
	if contact.ID > 0 {
		tx = tx.Where("id = ?", contact.ID)
	} else {
		contact.CreatedAt = currTime
	}

	if err := tx.Save(contact).Error; err != nil {
		return err
	}
	return nil
}
