package roles

import (
	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"

	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/logger"
	"github.com/costa92/go-web/pkg/util"
)

type RequestDetail struct {
	Id int `query:"id" binding:"gte=1" form:"id"`
}
type DetailResponse struct {
	*model.Role
	Menus []int `json:"menus"`
}

func (api *RoleController) Detail(ctx *gin.Context) {
	var req RequestDetail
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
	roleMenuModel := model.NewRoleMenuModel(ctx, api.MysqlStorage)
	menuIds, err := roleMenuModel.GetRolesByRoleId(req.Id)
	if err != nil {
		logger.Errorw("RoleController Detail NewRoleMenuModel.GetRolesByRoleId", "err", err)
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), nil)
		return
	}

	util.WriteResponse(ctx, nil, DetailResponse{Role: &role, Menus: menuIds})
}

func (api *RoleController) Del(ctx *gin.Context) {
}
