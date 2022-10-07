package roles

import (
	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"

	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/util"
)

type UpdateRequest struct {
	CreateRequest
	Id int `form:"id" binding:"gte=1" json:"id"`
}

func (api *RoleController) Update(ctx *gin.Context) {
	var req UpdateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrValidation, err.Error()), "参数错误")
		return
	}
	var role model.Role
	tx := api.MysqlStorage.Model(&model.Role{})
	if err := tx.Where("id = ?", req.Id).First(&role).Error; err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "查询数据错误")
		return
	}
	role.Name = req.Name
	role.Remark = req.Remark
	if err := tx.Save(&role).Error; err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "保存数据错误")
		return
	}
	util.WriteResponse(ctx, nil, "success")
}
