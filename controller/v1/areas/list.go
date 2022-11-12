package areas

import (
	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"

	"github.com/costa92/go-web/internal/validation"
	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/util"
)

func (c *AreaController) GetProvinces(ctx *gin.Context) {
	areaModel := model.NewAreaModel(ctx, c.MysqlStorage)
	areas, err := areaModel.QueryProvinces()
	if err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()),
			"省份列表数据错误")
		return
	}
	util.WriteSuccessResponse(ctx, areas)
}

type GetAreasByPidRequest struct {
	ParentCode int64 `query:"parentCode" binding:"gte=0" form:"parentCode"`
}

func (c *AreaController) GetAreasByPid(ctx *gin.Context) {
	var req GetAreasByPidRequest
	if err := ctx.ShouldBind(&req); err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrBind, err.Error()), nil)
		return
	}
	validate := validation.DefaultValidator{}
	err := validate.ValidateStruct(req)
	if err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrValidation, err.Error()), nil)
		return
	}
	areaModel := model.NewAreaModel(ctx, c.MysqlStorage)
	areas, err := areaModel.QueryListByPid(req.ParentCode)
	if err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()),
			"地区列表数据错误")
		return
	}
	util.WriteSuccessResponse(ctx, areas)
}
