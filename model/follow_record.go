package model

import (
	"context"
	"github.com/costa92/go-web/pkg/meta"
	"gorm.io/gorm"
	"time"
)

const TableNameFollowRecord = "follow_record"

// FollowRecord mapped from table <follow_record>
type FollowRecord struct {
	RecordID     int32  `gorm:"column:record_id;primaryKey;autoIncrement:true" json:"record_id"` // 主键
	EnterpriseID int32  `gorm:"column:enterprise_id;not null" json:"enterprise_id"`              // 企业编号
	Message      string `gorm:"column:message;not null" json:"message"`                          // 提交内容
	Status       int32  `gorm:"column:status;not null" json:"status"`                            // 跟进的状态
	UserID       int32  `gorm:"column:user_id;not null" json:"user_id"`                          // 跟进人
	ContactID    int32  `gorm:"column:contact_id;not null" json:"contact_id"`                    // 企业联系人
	UpdatedAt    int64  `gorm:"column:updated_at;not null" json:"updated_at"`                    // 修改时间
	UpdatedBy    int32  `gorm:"column:updated_by;not null" json:"updated_by"`                    // 修改者
	CreatedAt    int64  `gorm:"column:created_at;not null" json:"created_at"`                    // 添加时间
	CreatedBy    int32  `gorm:"column:created_by;not null" json:"created_by"`                    // 添加者
}

// TableName FollowRecord's table name
func (*FollowRecord) TableName() string {
	return TableNameFollowRecord
}

type FollowRecordList struct {
	meta.ListMeta `json:",inline"`
	Items         []*FollowRecord `json:"items"`
}
type FollowRecordModel struct {
	DB *gorm.DB
}

func NewFollowRecordModel(ctx context.Context, db *gorm.DB) *FollowRecordModel {
	return &FollowRecordModel{
		DB: db.Model(&FollowRecord{}).WithContext(ctx),
	}
}

func (m *FollowRecordModel) FindByEnterpriseId(enterpriseId int64) ([]*FollowRecord, error) {
	var records []*FollowRecord
	tx := m.DB.Where("enterprise_id = ?", enterpriseId).Order("record_id desc")
	if err := tx.Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

func (m *FollowRecordModel) Save(ctx context.Context, record *FollowRecord) error {
	tx := m.DB
	currTime := time.Now().Unix()
	record.UpdatedAt = currTime
	if record.RecordID > 0 {
		tx = tx.Where("id = ?", record.RecordID)
	} else {
		record.CreatedAt = currTime
	}

	if err := tx.Save(record).Error; err != nil {
		return err
	}
	return nil
}
