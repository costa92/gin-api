package follow

import (
	"context"
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

type GetFollowListRequest struct {
	model.PageRequest
}

type FollowItem struct {
	*model.Distribute
	*model.Enterprise
	DistributeId   int32  `json:"distribute_id"`
	AreaName       string `json:"area_name"`
	CreatedTime    string `json:"created_time"`
	UpdatedTime    string `json:"updated_time"`
	EnterpriseType string `json:"enterprise_type"`
}

type GetFollowListResponse struct {
	Items         []*FollowItem `json:"items"`
	meta.ListMeta `json:",inline"`
}

func (c *FollowController) GetFollowList(ctx *gin.Context) {
	var req GetFollowListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		util.WriteErrResponse(ctx, code.ErrBind, err, nil)
		return
	}
	ret := &model.DistributeList{}
	distributeModel := model.NewDistributeModel(ctx, c.MysqlStorage)
	tx := distributeModel.DB.Where("status = ?", 1)

	if err := tx.Count(&ret.TotalCount).Error; err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "查询数据错误")
		return
	}
	if err := tx.Scopes(gormutil.Paginate(int(req.Page), int(req.PageSize))).
		Order("id desc").Find(&ret.Items).Error; err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "查询数据错误")
		return
	}
	items := make([]*FollowItem, 0)
	if len(ret.Items) > 0 {
		typeAll, err := c.getEnterpriseType(ctx)
		if err != nil {
			util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "查询数据错误")
			return
		}

		var enterpriseIds []uint
		for _, item := range ret.Items {
			enterpriseIds = append(enterpriseIds, item.EnterpriseID)
		}

		enterpriseModel := model.NewEnterpriseModel(ctx, c.MysqlStorage)
		enterprises, err := enterpriseModel.FindByIds(enterpriseIds)
		if err != nil {
			util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "查询数据错误")
			return
		}

		enterpriseMap := make(map[uint]*model.Enterprise)
		var areaIds []int
		for _, enterprise := range enterprises {
			enterpriseMap[enterprise.Id] = enterprise
			areaIds = append(areaIds, enterprise.ProvinceId, enterprise.CityId, enterprise.CountyId)
		}

		areas, err := c.getAreas(ctx, areaIds)
		if err != nil {
			util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "查询数据错误")
			return
		}

		for _, item := range ret.Items {
			enterpriseItem := enterpriseMap[item.EnterpriseID]
			followItem := &FollowItem{
				Distribute:   item,
				DistributeId: item.ID,
				Enterprise:   enterpriseMap[item.EnterpriseID],
			}
			var areaName string
			var enterpriseType int
			if enterpriseItem != nil {
				areaName = fmt.Sprintf("%s%s%s", areas[enterpriseItem.ProvinceId], areas[enterpriseItem.CityId], areas[enterpriseItem.CountyId])
				enterpriseType = enterpriseMap[item.EnterpriseID].Type
				if enterpriseItem.CreatedAt > 0 {
					followItem.CreatedTime = time.Unix(item.CreatedAt, 0).Format(middleware.TimeFieldFormat)
				}
				if enterpriseItem.UpdateAt > 0 {
					followItem.UpdatedTime = time.Unix(item.UpdatedAt, 0).Format(middleware.TimeFieldFormat)
				}
				followItem.EnterpriseType = typeAll[enterpriseType]
				followItem.AreaName = areaName
			}

			items = append(items, followItem)
		}
	}
	util.WriteSuccessResponse(ctx, GetFollowListResponse{Items: items, ListMeta: ret.ListMeta})
}

func (c *FollowController) getEnterpriseType(ctx context.Context) (map[int]string, error) {
	enterpriseTypeModel := model.NewEnterpriseTypeModel(ctx, c.MysqlStorage)
	enterpriseTypes, err := enterpriseTypeModel.QueryList("")
	if err != nil {
		return nil, err
	}
	enterpriseTypeNames := make(map[int]string)
	for _, enterpriseType := range enterpriseTypes {
		enterpriseTypeNames[int(enterpriseType.ID)] = enterpriseType.Name
	}
	return enterpriseTypeNames, err
}

func (c *FollowController) getAreas(ctx context.Context, areaIds []int) (map[int]string, error) {
	areaModel := model.NewAreaModel(ctx, c.MysqlStorage)
	areas, err := areaModel.QueryListByIds(areaIds)
	if err != nil {
		return nil, err
	}
	areaNames := make(map[int]string, 0)
	for _, area := range areas {
		areaNames[int(area.ID)] = area.Name
	}
	return areaNames, err
}
