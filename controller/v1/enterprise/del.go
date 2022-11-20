package enterprise

import (
	"github.com/costa92/errors"
	"github.com/costa92/go-web/internal/validation"
	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/util"
	"github.com/gin-gonic/gin"
)

type DeletedRequest struct {
	Id int `query:"id" binding:"gte=1" form:"id"`
}

func (c *EnterpriseController) Deleted(ctx *gin.Context) {
	var req DeletedRequest
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
	enterpriseModel := model.NewEnterpriseModel(ctx, c.MysqlStorage)
	enterprise, err := enterpriseModel.FirstById(req.Id)
	if err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "查询企业数据错误")
		return
	}

	if err := enterpriseModel.Deleted(enterprise.Id); err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "删除企业数据失败")
		return
	}
	util.WriteSuccessResponse(ctx, "success")
}
