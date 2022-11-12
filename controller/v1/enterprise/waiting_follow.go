package enterprise

import (
	"fmt"
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

type GetWaitingFollowedListRequest struct {
	model.PageRequest
	Type int    `json:"type" query:"type" form:"type"`
	Name string `json:"name" query:"name" form:"name"`
}

type GetWaitingFollowedListResponse struct {
	Items         []*EnterpriseItem `json:"items"`
	meta.ListMeta `json:",inline"`
}

func (c *EnterpriseController) GetWaitingFollowedList(ctx *gin.Context) {
	var req GetWaitingFollowedListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		util.WriteErrResponse(ctx, code.ErrBind, err, nil)
		return
	}
	ret := &model.EnterpriseList{}
	enterpriseModel := model.NewEnterpriseModel(ctx, c.MysqlStorage)
	tx := enterpriseModel.DB.Where("status =? ", 1)
	if req.Type > 0 {
		tx = tx.Where("type =?", req.Type)
	}
	if req.Name != "" {
		tx = tx.Where("name like ?", "%"+req.Name+"%")
	}
	if err := tx.Count(&ret.TotalCount).Error; err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "查询数据错误")
		return
	}

	if err := tx.Scopes(gormutil.Paginate(int(req.Page), int(req.PageSize))).
		Order("id desc").Find(&ret.Items).Error; err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "查询数据错误")
		return
	}
	items := make([]*EnterpriseItem, 0)
	if len(ret.Items) > 0 {
		var areaIds []int
		for _, item := range ret.Items {
			areaIds = append(areaIds, item.ProvinceId, item.CityId, item.CountyId)
		}
		areaModel := model.NewAreaModel(ctx, c.MysqlStorage)
		areas, err := areaModel.QueryListByIds(areaIds)
		if err != nil {
			util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "查询数据错误")
		}

		areaNames := make(map[int]string, 0)
		for _, area := range areas {
			areaNames[int(area.ID)] = area.Name
		}

		enterpriseTypeModel := model.NewEnterpriseTypeModel(ctx, c.MysqlStorage)
		enterpriseTypes, err := enterpriseTypeModel.QueryList("")
		if err != nil {
			util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "查询数据错误")
		}
		enterpriseTypeNames := make(map[int]string)
		for _, enterpriseType := range enterpriseTypes {
			enterpriseTypeNames[int(enterpriseType.ID)] = enterpriseType.Name
		}

		for _, item := range ret.Items {
			areaName := fmt.Sprintf("%s%s%s", areaNames[item.ProvinceId], areaNames[item.CityId], areaNames[item.CountyId])
			enterpriseItem := &EnterpriseItem{
				Enterprise: item,
				AreaName:   areaName,
			}
			if enterpriseItem.CreatedAt > 0 {
				enterpriseItem.CreatedTime = time.Unix(enterpriseItem.CreatedAt, 0).Format(middleware.TimeFieldFormat)
			}
			if enterpriseItem.UpdateAt > 0 {
				enterpriseItem.UpdatedTime = time.Unix(enterpriseItem.UpdateAt, 0).Format(middleware.TimeFieldFormat)
			}
			if enterpriseItem.Type > 0 {
				enterpriseItem.EnterpriseType = enterpriseTypeNames[enterpriseItem.Type]
			}
			items = append(items, enterpriseItem)
		}
	}

	util.WriteSuccessResponse(ctx, GetWaitingFollowedListResponse{Items: items, ListMeta: ret.ListMeta})
	return
}
