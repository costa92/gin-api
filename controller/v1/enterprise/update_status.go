package enterprise

import (
	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"

	"github.com/costa92/go-web/internal/validation"
	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/util"
)

func (c *EnterpriseController) UpdateStatus(ctx *gin.Context) {
	var req UpdateStatusRequest
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

	enterprise.Status = req.Status
	if err := enterpriseModel.Save(enterprise); err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "更新企业状态错误")
		return
	}

	util.WriteSuccessResponse(ctx, enterprise)
}
