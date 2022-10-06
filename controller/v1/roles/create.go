package roles

import (
	"github.com/costa92/errors"
	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/util"
	"github.com/gin-gonic/gin"
)

type CreateRequest struct {
	Name   string `form:"Name" binding:"required" json:"name"`
	Remark string `form:"remark" binding:"required" json:"remark"`
}

func (api *RoleController) Create(ctx *gin.Context) {
	var rep CreateRequest
	if err := ctx.ShouldBind(&rep); err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrValidation, err.Error()), "参数错误")
		return
	}
	role := &model.Role{
		Name:   rep.Name,
		Remark: rep.Remark,
	}
	if err := api.MysqlStorage.Model(&model.Role{}).Create(role).Error; err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "添加角色错误")
		return
	}
	util.WriteResponse(ctx, nil, "success")
}
