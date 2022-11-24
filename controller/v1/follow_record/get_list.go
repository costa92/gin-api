package follow_record

import (
	"time"

	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"

	"github.com/costa92/go-web/internal/middleware"
	"github.com/costa92/go-web/internal/validation"
	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/meta"
	"github.com/costa92/go-web/pkg/util"
	"github.com/costa92/go-web/pkg/util/gormutil"
)

type FollowRecordItem struct {
	*model.FollowRecord
	CreatedTime string `json:"created_time"`
	UpdatedTime string `json:"updated_time"`
}

type FollowRecordsResponse struct {
	Items         []*FollowRecordItem `json:"items"`
	meta.ListMeta `json:",inline"`
}

func (m *FollowRecordController) FollowRecords(ctx *gin.Context) {
	var req FollowRecordRequest
	if err := ctx.ShouldBind(&req); err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrBind, err.Error()), nil)
		return
	}
	validate := validation.DefaultValidator{}
	err := validate.ValidateStruct(req)
	if err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrValidation, err.Error()), nil)
		return
	}

	ret := &model.FollowRecordList{}

	tx := m.MysqlStorage.WithContext(ctx).Model(&model.FollowRecord{})
	tx = tx.Where("enterprise_id = ?", req.EnterpriseId)

	if err := tx.Count(&ret.TotalCount).Error; err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "查询数据错误")
		return
	}

	if err := tx.Scopes(gormutil.Paginate(int(req.Page), int(req.PageSize))).
		Order("record_id desc").Find(&ret.Items).Error; err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "查询数据错误")
		return
	}

	items := make([]*FollowRecordItem, 0)
	if len(ret.Items) > 0 {
		for _, record := range ret.Items {
			items = append(items, &FollowRecordItem{
				FollowRecord: record,
				CreatedTime:  time.Unix(record.CreatedAt, 0).Format(middleware.TimeFieldFormat),
			})
		}
	}
	util.WriteSuccessResponse(ctx, FollowRecordsResponse{Items: items, ListMeta: ret.ListMeta})
}
