package enterprise_type

import (
	"github.com/gin-gonic/gin"

	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/util"
)

func (c *EnterpriseTypeController) Create(ctx *gin.Context) {
	var req CreateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		util.WriteErrResponse(ctx, code.ErrBind, err, "参数错误")
		return
	}
	eTypeModel := model.NewEnterpriseTypeModel(ctx, c.MysqlStorage)
	eType := &model.EnterpriseType{}
	c.SaveParams(eType, &req)
	if err := eTypeModel.Save(eType); err != nil {
		util.WriteResponse(ctx, err, "添加企业类型错误")
		return
	}
	util.WriteResponse(ctx, nil, "success")
}
