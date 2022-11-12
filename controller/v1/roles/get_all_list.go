package roles

import (
	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"

	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/util"
)

func (api *RoleController) GetAllList(ctx *gin.Context) {
	roles := make([]*model.Role, 0)
	tx := api.MysqlStorage.Model(&model.Role{})
	if err := tx.Where("status =?", 1).Find(&roles).Error; err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "查询数据库错误")
		return
	}
	util.WriteSuccessResponse(ctx, roles)
}
