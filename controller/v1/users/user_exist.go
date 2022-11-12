package users

import (
	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/logger"
	"github.com/costa92/go-web/pkg/util"
)

type GetUserExitRequest struct {
	Username string `json:"username,omitempty"  form:"username"`
}

func (u *UserController) PostUserAccountExit(ctx *gin.Context) {
	var req GetUserExitRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrBind, err.Error()), "参数错误")
		return
	}
	userModel := model.NewUserModel(ctx, u.MysqlStorage)
	user, err := userModel.FirstByName(req.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Errorw("UserController Get FirstByName", "err", err)
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), nil)
		return
	}
	if user != nil {
		util.WriteSuccessResponse(ctx, false)
		return
	}
	util.WriteSuccessResponse(ctx, true)
}
