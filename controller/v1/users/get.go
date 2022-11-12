package users

import (
	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"

	"github.com/costa92/go-web/internal/validation"
	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/logger"
	"github.com/costa92/go-web/pkg/util"
)

type DetailRequest struct {
	*model.User
	Role []int `json:"role"`
}

func (u *UserController) Get(ctx *gin.Context) {
	var req GetUserRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrBind, err.Error()), nil)
		return
	}
	validate := validation.DefaultValidator{}
	err := validate.ValidateStruct(req)
	if err != nil {
		logger.Errorf("get", "err", err, "req", req)
		util.WriteResponse(ctx, errors.WithCode(code.ErrValidation, err.Error()), nil)
		return
	}
	userModel := model.NewUserModel(ctx, u.MysqlStorage)
	user, err := userModel.FirstByUid(req.Id)
	if err != nil {
		logger.Errorw("UserController Get FirstByUid", "err", err)
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), nil)
		return
	}
	userRoleModel := model.NewUserRoleModel(ctx, u.MysqlStorage)
	roleIds, err := userRoleModel.GetRolesByUserId(user.ID)
	if err != nil {
		logger.Errorw("UserController NewUserRoleModel.GetRolesByUserId", "err", err)
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), nil)
		return
	}
	util.WriteSuccessResponse(ctx, DetailRequest{User: user, Role: roleIds})
}
