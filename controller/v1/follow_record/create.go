package follow_record

import (
	"github.com/costa92/errors"
	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/util"
	"github.com/gin-gonic/gin"
)

func (m *FollowRecordController) Create(ctx *gin.Context) {
	var req CreateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrBind, err.Error()), nil)
		return
	}
	record := &model.FollowRecord{}
	record = m.SaveParams(record, &req)
	recordModel := model.NewFollowRecordModel(ctx, m.MysqlStorage)
	if err := recordModel.Save(ctx, record); err != nil {
		util.WriteResponse(ctx, err, "添加记录错误")
		return
	}
	util.WriteSuccessResponse(ctx, "success")
}
