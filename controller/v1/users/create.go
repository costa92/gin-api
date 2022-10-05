package users

import (
	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/costa92/go-web/internal/validation"
	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/logger"
	"github.com/costa92/go-web/pkg/util"
	"github.com/costa92/go-web/pkg/util/auth"
)

var defaultPassword = "123456"

func (u *UserController) Create(ctx *gin.Context) {
	var req UserCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrBind, err.Error()), nil)
		return
	}
	validate := validation.DefaultValidator{}
	err := validate.ValidateStruct(req)
	if err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrValidation, err.Error()), nil)
		return
	}
	userModel := model.NewUserModel(ctx, u.MysqlStorage)
	userInfo, err := userModel.FirstByName(req.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), nil)
		return
	}
	if userInfo != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrUserAlreadyExist, "用户已经存在"), nil)
		return
	}

	user := &model.User{}
	password, _ := auth.Encrypt(defaultPassword)
	u.saveParams(user, &req)
	user.Password = password
	if err := userModel.Save(user); err != nil {
		logger.Errorf("userController Create Save failed", "err", err)
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), nil)
		return
	}
	util.WriteSuccessResponse(ctx, "success")
}
