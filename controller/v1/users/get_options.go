package users

import (
	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"

	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/util"
)

type GetUserOptionsItem struct {
	Value int32  `json:"value"`
	Label string `json:"label"`
}

func (u *UserController) GetOptions(ctx *gin.Context) {
	userModel := model.NewUserModel(ctx, u.MysqlStorage)
	users, err := userModel.QueryList("")
	if err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "查询类型数据错误")
		return
	}
	items := make([]*GetUserOptionsItem, 0)
	if len(users) > 0 {
		for _, item := range users {
			items = append(items, &GetUserOptionsItem{Value: int32(item.ID), Label: item.RealName})
		}
	}
	util.WriteSuccessResponse(ctx, items)
}
