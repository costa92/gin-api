package roles

import (
	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"

	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/meta"
	"github.com/costa92/go-web/pkg/util"
	"github.com/costa92/go-web/pkg/util/gormutil"
)

func (api *RoleController) Index(ctx *gin.Context) {
	var r meta.ListOptions
	ret := &model.RoleList{}
	if err := ctx.ShouldBindQuery(&r); err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrBind, err.Error()), nil)
		return
	}
	// 分页处理
	ol := gormutil.Unpointer(r.Offset, r.Limit)
	// 处理查询数据
	if err := api.MysqlStorage.
		Model(&model.Role{}).Scopes(gormutil.Paginate(ol.Offset, ol.Limit)).
		Order("id desc").
		Find(&ret.Items).
		Offset(-1).
		Limit(-1).
		Count(&ret.TotalCount).
		Error; err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "查询数据库错误")
		return
	}

	util.WriteSuccessResponse(ctx, ret)
}
