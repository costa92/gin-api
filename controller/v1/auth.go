package v1

import (
	"strconv"
	"strings"

	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/costa92/go-web/internal/middleware"
	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/logger"
	"github.com/costa92/go-web/pkg/util"
)

type AuthController struct {
	MysqlStorage *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{MysqlStorage: db}
}

type GetUserInfoResponse struct {
	*model.User
	UserId   string              `json:"userId"`
	RealName string              `json:"realName"`
	Avatar   string              `json:"avatar"`
	Desc     string              `json:"desc"`
	Token    string              `json:"token"`
	Roles    []map[string]string `json:"roles"`
}

func (a *AuthController) GetUserInfo(ctx *gin.Context) {
	author := strings.SplitN(ctx.Request.Header.Get("Authorization"), " ", 2)
	var token string
	if len(author) == 2 {
		token = author[1]
	}
	userModel := model.NewUserModel(ctx, a.MysqlStorage)
	username := ctx.GetString(middleware.UsernameKey)
	user, err := userModel.FirstByName(username)
	if err != nil {
		logger.Errorw("GetUserInfo FirstByName failed", middleware.UsernameKey, username, "err", err)
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "获取用户信息错误")
	}
	var roles []map[string]string
	roles = append(roles, map[string]string{"roleName": "test", "value": "123123"})
	resp := &GetUserInfoResponse{
		User:   user,
		UserId: strconv.Itoa(user.ID),
		Roles:  roles,
		Token:  token,
	}
	util.WriteResponse(ctx, nil, resp)
}
