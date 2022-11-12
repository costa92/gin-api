package roles

import (
	"github.com/costa92/errors"
	"github.com/costa92/go-web/internal/middleware"
	"github.com/gin-gonic/gin"
	"time"

	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/meta"
	"github.com/costa92/go-web/pkg/util"
	"github.com/costa92/go-web/pkg/util/gormutil"
)

type RoleItem struct {
	*model.Role
	CreatedTime string `json:"created_time"`
}

type IndexResponse struct {
	Items         []*RoleItem `json:"items"`
	meta.ListMeta `json:",inline"`
}

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
	items := make([]*RoleItem, 0)
	if len(ret.Items) > 0 {
		for _, item := range ret.Items {
			var createTime string
			if item.CreatedAt > 0 {
				createTime = time.Unix(item.CreatedAt, 0).Format(middleware.TimeFieldFormat)
			}
			items = append(items, &RoleItem{
				Role:        item,
				CreatedTime: createTime,
			})
		}
	}

	util.WriteSuccessResponse(ctx, IndexResponse{Items: items, ListMeta: ret.ListMeta})
}
