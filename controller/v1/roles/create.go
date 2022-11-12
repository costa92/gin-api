package roles

import (
	"context"
	"time"

	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"

	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/util"
)

type CreateRequest struct {
	Name   string `form:"name" binding:"required" json:"name"`
	Remark string `form:"remark" binding:"required" json:"remark"`
	Status int    `form:"status" binding:"gte=1"  json:"status"`
	Menus  []int  `form:"menus" binding:"required" json:"menus"`
}

func (api *RoleController) Create(ctx *gin.Context) {
	var rep CreateRequest
	if err := ctx.ShouldBind(&rep); err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrValidation, err.Error()), "参数错误")
		return
	}
	currUnix := time.Now().Unix()
	role := &model.Role{
		Name:      rep.Name,
		Remark:    rep.Remark,
		Status:    rep.Status,
		UpdatedAt: currUnix,
		CreatedAt: currUnix,
	}
	if err := api.MysqlStorage.Model(&model.Role{}).Debug().Create(role).Error; err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "添加角色错误")
		return
	}
	if err := api.saveRoleMenus(ctx, int(role.ID), rep.Menus); err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "添加角色菜单错误")
		return
	}
	util.WriteResponse(ctx, nil, "success")
}

func (api *RoleController) saveRoleMenus(ctx context.Context, roleId int, menuIds []int) error {
	roleMenuModel := model.NewRoleMenuModel(ctx, api.MysqlStorage)
	if err := roleMenuModel.DeletedByRoleId(roleId); err != nil {
		return err
	}
	roleMenus := make([]model.RoleMenu, 0)
	for _, id := range menuIds {
		roleMenus = append(roleMenus, model.RoleMenu{
			MenuId: uint(id),
			RoleId: uint(roleId),
		})
	}
	return roleMenuModel.DB.Create(&roleMenus).Error
}
