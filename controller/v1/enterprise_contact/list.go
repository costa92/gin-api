package enterprise_contact

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

type GeListRequest struct {
	model.PageRequest
	Name   string `json:"name,omitempty" form:"name" query:"name"`
	Mobile string `json:"mobile,omitempty" form:"mobile" query:"mobile"`
}
type GeListItem struct {
	*model.EnterpriseContact
	EnterpriseName string `json:"enterprise_name"`
	CreatedTime    string `json:"created_time"`
	UpdatedTime    string `json:"updated_time"`
}
type GetListResponse struct {
	Items         []*GeListItem `json:"items"`
	meta.ListMeta `json:",inline"`
}

func (c *EnterpriseContactController) GeList(ctx *gin.Context) {
	var req GeListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		util.WriteErrResponse(ctx, code.ErrBind, err, nil)
		return
	}

	ret := &model.EnterpriseContactList{}
	contactModel := model.NewEnterpriseContactModel(ctx, c.MysqlStorage)
	tx := contactModel.DB.Where("status = ?", model.ContactStatusNormal)

	if req.Name != "" {
		tx = tx.Where("name like ?", "%"+req.Name+"%")
	}
	if req.Mobile != "" {
		tx = tx.Where("mobile like ?", "%"+req.Mobile+"%")
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

	items := make([]*GeListItem, 0)
	if len(ret.Items) > 0 {
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
		enterpriseArr := make(map[uint]string)
		for _, enterprise := range enterprises {
			enterpriseArr[enterprise.Id] = enterprise.Name
		}

		for _, item := range ret.Items {
			items = append(items, &GeListItem{
				EnterpriseContact: item,
				CreatedTime:       time.Unix(item.CreatedAt, 0).Format(middleware.TimeFieldFormat),
				UpdatedTime:       time.Unix(item.UpdatedAt, 0).Format(middleware.TimeFieldFormat),
				EnterpriseName:    enterpriseArr[item.EnterpriseID],
			})
		}
	}

	util.WriteSuccessResponse(ctx, GetListResponse{
		Items:    items,
		ListMeta: ret.ListMeta,
	})
}
