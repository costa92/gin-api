package users

import (
	"gorm.io/gorm"

	"github.com/costa92/go-web/model"
)

type UserController struct {
	MysqlStorage *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{MysqlStorage: db}
}

type UserCreateRequest struct {
	Username string `json:"username" validate:"required"`
	Nickname string `json:"nickname" validate:"required"`
	Mobile   int64  `json:"mobile" validate:"required"`
	RealName string `json:"real_name" validate:"required"`
	Status   int    `json:"status" validate:"oneof=1 2"`
}

func (u *UserController) saveParams(user *model.User, req *UserCreateRequest) {
	user.Username = req.Username
	user.Nickname = req.Nickname
	user.Mobile = req.Mobile
	user.Status = req.Status
	user.RealName = req.RealName
}

type UserUpdateRequest struct {
	GetUserRequest
	UserCreateRequest
}

type UserUpdateStatusRequest struct {
	GetUserRequest
	Status int `json:"status" validate:"oneof=1 2"`
}

type GetUserRequest struct {
	Id int `json:"id" form:"id" validate:"required"`
}
