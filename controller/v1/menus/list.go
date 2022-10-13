package menus

import (
	"github.com/costa92/errors"
	"github.com/costa92/go-web/internal/middleware"
	"github.com/gin-gonic/gin"
	"time"

	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/util"
)

type MenuListRequest struct {
	model.PageRequest
}

func (m *MenuController) List(ctx *gin.Context) {
	var req MenuListRequest
	ret := &model.MenuList{}

	if err := ctx.ShouldBindQuery(&req); err != nil {
		util.WriteErrResponse(ctx, code.ErrBind, err, nil)
		return
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrBind, err.Error()), nil)
		return
	}
	// 处理查询数据
	if err := m.MysqlStorage.
		Model(&model.Menu{}).
		Order("menu_sort asc").
		Find(&ret.Items).
		Error; err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "查询数据库错误")
		return
	}
	trees := listMenuTree(ret.Items, 0)
	util.WriteSuccessResponse(ctx, trees)

}

type ListMenuTreeItem struct {
	*model.Menu
	FormattedCreatedAt string             `json:"formatted_created_at"`
	FormattedUpdatedAt string             `json:"formatted_updated_at"`
	Children           []ListMenuTreeItem `json:"children"`
}

type GetMenuTreeResponse struct {
	Items []ListMenuTreeItem `json:"items"`
}

func listMenuTree(menus []*model.Menu, parentID uint) []ListMenuTreeItem {
	tree := make([]ListMenuTreeItem, 0)
	for _, menu := range menus {
		if parentID != menu.ParentID {
			continue
		}
		item := ListMenuTreeItem{
			Menu:               menu,
			FormattedCreatedAt: time.Unix(menu.CreatedAt, 0).Format(middleware.TimeFieldFormat),
			FormattedUpdatedAt: time.Unix(menu.UpdatedAt, 0).Format(middleware.TimeFieldFormat),
			Children:           []ListMenuTreeItem{},
		}
		item.Children = listMenuTree(menus, menu.ID)
		tree = append(tree, item)
	}
	return tree
}
