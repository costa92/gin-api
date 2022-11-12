package enterprise

import (
	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"

	"github.com/costa92/go-web/internal/validation"
	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/util"
)

func (c *EnterpriseController) Update(ctx *gin.Context) {
	var req UpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
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

	c.saveParams(enterprise, &req.CreateRequest)
	if err := enterpriseModel.Save(enterprise); err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()),
			"企业修改类型数据错误")
		return
	}

	util.WriteSuccessResponse(ctx, "success")
}
