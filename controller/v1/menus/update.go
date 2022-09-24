package menus

import (
	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"

	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/util"
)

func (m *MenuController) Update(ctx *gin.Context) {
	var req MenuUpdateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrBind, err.Error()), "参数错误")
		return
	}
	menuModel := model.NewMenuModel(ctx, m.MysqlStorage)
	menu, err := menuModel.FirstByID(req.ID)
	if err != nil {
		util.WriteResponse(ctx, err, "添加菜单错误")
		return
	}
	menu = m.SaveParams(menu, &req.MenuCreateRequest)
	if err := menuModel.Save(ctx, menu); err != nil {
		util.WriteResponse(ctx, err, "添加菜单错误")
		return
	}
	util.WriteResponse(ctx, nil, "success")
}
