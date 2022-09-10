package auth

import (
	"github.com/costa92/errors"
	"github.com/costa92/go-web/internal/pkg/code"
	"github.com/costa92/go-web/internal/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type LoginRequest struct {
	UserName string `json:"user_name"  validate:"required" `
	Password string `json:"password"   validate:"required"`
}

func (a *Auth) Login(ctx *gin.Context) {
	var req LoginRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrBind, err.Error()), "参数错误")
		return
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrValidation, err.Error()), "参数验证错误")
		return
	}
	util.WriteResponse(ctx, nil, "ret")
}
