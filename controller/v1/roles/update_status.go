package roles

import (
	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/util"
	"github.com/gin-gonic/gin"
)

type UpdateStateRequest struct {
	Id     int `json:"id" binding:"gte=1" form:"id"`
	Status int `json:"status" binding:"required,oneof=1 2" form:"status"`
}

func (api *RoleController) UpdateState(ctx *gin.Context) {
	var req UpdateStateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		util.WriteErrResponse(ctx, code.ErrValidation, err, "参数错误")
		return
	}
	var role model.Role
	tx := api.MysqlStorage.Model(&model.Role{})
	if err := tx.Where("id = ?", req.Id).First(&role).Error; err != nil {
		util.WriteErrResponse(ctx, code.ErrDatabase, err, "查询数据错误")
		return
	}
	role.Status = req.Status
	if err := tx.Save(&role).Error; err != nil {
		util.WriteErrResponse(ctx, code.ErrDatabase, err, "修改角色状态错误")
		return
	}
	util.WriteSuccessResponse(ctx, "success")
}
