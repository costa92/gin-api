package menus

import (
	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"

	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/util"
)

func (m *MenuController) GetPermissionCode(ctx *gin.Context) {
	var menus []*model.Menu
	if err := m.MysqlStorage.Model(&model.Menu{}).Where("menu_type = ?", 3).Order("id desc").Find(&menus).Error; err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "查询数据库错误")
		return
	}
	util.WriteSuccessResponse(ctx, menus)
	return
}
