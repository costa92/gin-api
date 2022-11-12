package enterprise

import (
	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"

	"github.com/costa92/go-web/internal/validation"
	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/util"
)

func (c *EnterpriseController) Create(ctx *gin.Context) {
	var req CreateRequest
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
	enterprise := &model.Enterprise{}
	enterpriseModel := model.NewEnterpriseModel(ctx, c.MysqlStorage)
	c.saveParams(enterprise, &req)

	if err := enterpriseModel.Save(enterprise); err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()),
			"企业添加类型数据错误")
		return
	}
	util.WriteSuccessResponse(ctx, "success")
}
