package users

import (
	"context"

	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"

	"github.com/costa92/go-web/internal/validation"
	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/logger"
	"github.com/costa92/go-web/pkg/util"
)

func (u *UserController) Update(ctx *gin.Context) {
	var req UserUpdateRequest
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
	user, err := userModel.FirstByUid(req.Id)
	if err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), nil)
		return
	}
	tx := u.MysqlStorage.Begin()
	u.saveParams(user, &req.UserCreateRequest)
	if err := userModel.Save(user); err != nil {
		logger.Errorf("userController Update Save failed", "err", err)
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), nil)
		tx.Rollback()
		return
	}
	if err = u.saveUserRole(ctx, user.ID, req.Role); err != nil {
		logger.Errorf("userController saveUserRole  failed", "err", err)
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), nil)
		tx.Rollback()
		return
	}
	tx.Commit()
	util.WriteSuccessResponse(ctx, "success")
}

func (u *UserController) saveUserRole(ctx context.Context, userId int, roleIds []int) error {
	userRoleModel := model.NewUserRoleModel(ctx, u.MysqlStorage)
	if err := userRoleModel.DeletedByUserId(userId); err != nil {
		return err
	}
	userRoles := make([]model.UserRole, 0)
	for _, id := range roleIds {
		userRoles = append(userRoles, model.UserRole{
			UserID: userId,
			RoleId: id,
		})
	}
	return userRoleModel.DB.Create(&userRoles).Error
}
