package enterprise_type

import (
	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"

	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/util"
)

func (c *EnterpriseTypeController) Detail(ctx *gin.Context) {
	var req RequestDetail
	if err := ctx.ShouldBind(&req); err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrValidation, err.Error()), "参数错误")
		return
	}
	eTypeModel := model.NewEnterpriseTypeModel(ctx, c.MysqlStorage)
	eType, err := eTypeModel.FirstByID(req.ID)
	if err != nil {
		util.WriteResponse(ctx, err, "查询企业类型错误")
		return
	}
	util.WriteResponse(ctx, nil, eType)
}
