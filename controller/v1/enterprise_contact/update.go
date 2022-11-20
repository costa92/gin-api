package enterprise_contact

import (
	"github.com/costa92/errors"
	"github.com/costa92/go-web/internal/validation"
	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/util"
	"github.com/gin-gonic/gin"
)

func (c *EnterpriseContactController) Update(ctx *gin.Context) {
	var req UpdateRequest
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
	contactModel := model.NewEnterpriseContactModel(ctx, c.MysqlStorage)
	contact, err := contactModel.FirstById(req.Id)
	if err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), nil)
		return
	}
	c.updateParams(contact, &req)
	if err := contactModel.Save(contact); err != nil {
		util.WriteResponse(ctx, err, "添加企业联系人错误")
		return
	}
	util.WriteSuccessResponse(ctx, "success")
}