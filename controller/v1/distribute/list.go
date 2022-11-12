package distribute

import (
	"github.com/gin-gonic/gin"

	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/util"
)

type GetListRequest struct {
	model.PageRequest
}

func (c *DistributeController) GetList(ctx *gin.Context) {
	var req GetListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		util.WriteErrResponse(ctx, code.ErrBind, err, nil)
		return
	}
}
