package menus

import (
	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"

	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/util"
)

func (m *MenuController) Create(ctx *gin.Context) {
	var req MenuCreateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrBind, err.Error()), nil)
		return
	}
	menu := &model.Menu{}
	menu = m.SaveParams(menu, &req)
	menuModel := model.NewMenuModel(ctx, m.MysqlStorage)
	if err := menuModel.Save(ctx, menu); err != nil {
		util.WriteResponse(ctx, err, "添加菜单错误")
		return
	}
	util.WriteSuccessResponse(ctx, "success")
}
