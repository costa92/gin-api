package users

import (
	"time"

	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"

	"github.com/costa92/go-web/internal/middleware"
	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/meta"
	"github.com/costa92/go-web/pkg/util"
	"github.com/costa92/go-web/pkg/util/gormutil"
)

type UserItem struct {
	*model.User
	StatusDesc string   `json:"status_desc"`
	CreatedAt  string   `json:"created_at"`
	RoleName   []string `json:"role_name"`
	LastTime   string   `json:"last_time"`
}

type UserResponse struct {
	Items         []*UserItem `json:"items"`
	meta.ListMeta `json:",inline"`
}

func (u *UserController) Users(ctx *gin.Context) {
	var r meta.ListOptions
	ret := &model.UserList{}
	if err := ctx.ShouldBindQuery(&r); err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrBind, err.Error()), "参数错误")
		return
	}
	// 分页处理
	ol := gormutil.Unpointer(r.Offset, r.Limit)
	// 处理查询数据
	if err := u.MysqlStorage.
		Model(&model.User{}).Scopes(gormutil.Paginate(ol.Offset, ol.Limit)).
		Order("id desc").
		Find(&ret.Items).
		Offset(-1).
		Limit(-1).
		Count(&ret.TotalCount).
		Error; err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "查询数据库错误")
		return
	}

	items := make([]*UserItem, 0)
	if len(ret.Items) > 0 {
		var userIds []int64
		for _, item := range ret.Items {
			userIds = append(userIds, int64(item.ID))
		}

		userRoleModel := model.NewUserRoleModel(ctx, u.MysqlStorage)
		userRole, err := userRoleModel.FindByUserIds(userIds)
		if err != nil {
			util.WriteResponse(ctx, errors.WithCode(code.ErrRoleUserNotFound, err.Error()), "查询用户关系错误")
			return
		}

		// 处理角色
		var roleIds []int64
		rolesRoleMap := make(map[int][]int64, 0)
		for _, role := range userRole {
			roleIds = append(roleIds, int64(role.RoleId))
			if ok := rolesRoleMap[role.UserID]; ok != nil {
				rolesRoleMap[role.UserID] = append(rolesRoleMap[role.UserID], int64(role.RoleId))
			} else {
				rolesRoleMap[role.UserID] = []int64{int64(role.RoleId)}
			}
		}

		roleModel := model.NewRoleModel(ctx, u.MysqlStorage)
		roles, err := roleModel.FindByRoleIds(roleIds)
		if err != nil {
			util.WriteResponse(ctx, errors.WithCode(code.ErrRoleNotFound, err.Error()), "查询用户关系错误")
			return
		}
		rolesArr := make(map[int64]string, 0)
		for _, role := range roles {
			rolesArr[role.ID] = role.Name
		}

		for _, user := range ret.Items {
			userRolesId := rolesRoleMap[user.ID]
			var roleName []string
			if len(userRolesId) > 0 {
				for _, roleId := range userRolesId {
					roleName = append(roleName, rolesArr[roleId])
				}
			}
			item := &UserItem{
				User:       user,
				StatusDesc: model.UserStatusOptionDesc.Option(user.Status),
				CreatedAt:  time.Unix(user.CreatedAt, 0).Format(middleware.TimeFieldFormat),
				RoleName:   roleName,
			}
			if user.LastTime > 0 {
				item.LastTime = time.Unix(user.LastTime, 0).Format(middleware.TimeFieldFormat)
			}
			items = append(items, item)
		}
	}
	util.WriteResponse(ctx, nil, UserResponse{Items: items, ListMeta: ret.ListMeta})
}
