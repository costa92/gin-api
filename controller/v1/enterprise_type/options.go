package enterprise_type

import (
	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"

	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/util"
)

type GetOptionsItem struct {
	Value int32  `json:"value"`
	Label string `json:"label"`
}

func (c *EnterpriseTypeController) GetOptions(ctx *gin.Context) {
	enterpriseTypeModel := model.NewEnterpriseTypeModel(ctx, c.MysqlStorage)
	types, err := enterpriseTypeModel.QueryList("")
	if err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "查询类型数据错误")
		return
	}
	items := make([]*GetOptionsItem, 0)
	if len(types) > 0 {
		for _, item := range types {
			items = append(items, &GetOptionsItem{Value: item.ID, Label: item.Name})
		}
	}
	util.WriteSuccessResponse(ctx, items)
}
