package enterprise_type

import (
	"github.com/gin-gonic/gin"

	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/util"
)

func (c *EnterpriseTypeController) Update(ctx *gin.Context) {
	var req UpdateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		util.WriteErrResponse(ctx, code.ErrBind, err, "参数错误")
		return
	}
	eTypeModel := model.NewEnterpriseTypeModel(ctx, c.MysqlStorage)
	eType, err := eTypeModel.FirstByID(req.ID)
	if err != nil {
		util.WriteResponse(ctx, err, "查询企业类型错误")
		return
	}
	c.SaveParams(eType, &req.CreateRequest)
	if err := eTypeModel.Save(eType); err != nil {
		util.WriteResponse(ctx, err, "添加企业类型错误")
		return
	}
	util.WriteResponse(ctx, nil, "success")
}
