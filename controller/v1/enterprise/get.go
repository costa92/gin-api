package enterprise

import (
	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"

	"github.com/costa92/go-web/internal/validation"
	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/util"
)

func (c *EnterpriseController) Detail(ctx *gin.Context) {
	var req DetailRequest
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
	enterpriseModel := model.NewEnterpriseModel(ctx, c.MysqlStorage)
	enterprise, err := enterpriseModel.FirstById(req.Id)
	if err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "查询企业数据错误")
		return
	}

	eTypeModel := model.NewEnterpriseTypeModel(ctx, c.MysqlStorage)
	eType, err := eTypeModel.FirstByID(enterprise.Type)
	if err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "查询企业类型数据错误")
		return
	}

	detail := DetailResponse{
		Enterprise: enterprise,
		TypeName:   eType.Name,
	}
	detail.AreaId = append(detail.AreaId, enterprise.ProvinceId, enterprise.CityId, enterprise.CountyId)

	areaModel := model.NewAreaModel(ctx, c.MysqlStorage)
	areas, err := areaModel.QueryListByIds(detail.AreaId)
	if err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "查询地区错误")
		return
	}
	for _, area := range areas {
		detail.AreaName = append(detail.AreaName, area.Name)
	}
	util.WriteSuccessResponse(ctx, detail)
}
