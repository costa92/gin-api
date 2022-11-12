package enterprise_type

import (
	"time"

	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"

	"github.com/costa92/go-web/internal/middleware"
	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/meta"
	"github.com/costa92/go-web/pkg/util"
	"github.com/costa92/go-web/pkg/util/gormutil"
)

type EnterpriseTypeItem struct {
	*model.EnterpriseType
	CreatedTime string `json:"created_time"`
	UpdatedTime string `json:"updated_time"`
}
type GetListResponse struct {
	Items         []*EnterpriseTypeItem `json:"items"`
	meta.ListMeta `json:",inline"`
}

func (c *EnterpriseTypeController) Index(ctx *gin.Context) {
	var r meta.ListOptions
	ret := &model.EnterpriseTypeList{}
	if err := ctx.ShouldBindQuery(&r); err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrBind, err.Error()), nil)
		return
	}
	// 分页处理
	ol := gormutil.Unpointer(r.Offset, r.Limit)
	// 处理查询数据
	if err := c.MysqlStorage.
		Model(&model.EnterpriseType{}).Scopes(gormutil.Paginate(ol.Offset, ol.Limit)).
		Order("id desc").
		Find(&ret.Items).
		Offset(-1).
		Limit(-1).
		Count(&ret.TotalCount).
		Error; err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "查询数据库错误")
		return
	}
	items := make([]*EnterpriseTypeItem, 0)
	if len(ret.Items) > 0 {
		for _, eType := range ret.Items {
			item := &EnterpriseTypeItem{
				EnterpriseType: eType,
			}
			if eType.CreatedAt > 0 {
				item.CreatedTime = time.Unix(eType.CreatedAt, 0).Format(middleware.TimeFieldFormat)
			}
			if eType.UpdatedAt > 0 {
				item.UpdatedTime = time.Unix(eType.UpdatedAt, 0).Format(middleware.TimeFieldFormat)
			}
			items = append(items, item)
		}
	}
	util.WriteSuccessResponse(ctx, GetListResponse{Items: items, ListMeta: ret.ListMeta})
}
