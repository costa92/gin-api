package follow_record

import (
	"github.com/costa92/go-web/model"
	"gorm.io/gorm"
)

type FollowRecordController struct {
	MysqlStorage *gorm.DB
}

func NewFollowRecordController(db *gorm.DB) *FollowRecordController {
	return &FollowRecordController{
		MysqlStorage: db,
	}
}

type CreateRequest struct {
	ContactId    int32  `json:"contact_id"`
	EnterpriseId int32  `json:"enterprise_id"`
	Message      string `json:"message"`
}

func (m *FollowRecordController) SaveParams(record *model.FollowRecord, req *CreateRequest) *model.FollowRecord {
	record.EnterpriseID = req.EnterpriseId
	record.Message = req.Message
	record.ContactID = req.ContactId
	record.Status = 2
	return record
}
