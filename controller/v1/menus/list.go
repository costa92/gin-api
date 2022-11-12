package menus

import (
	"time"

	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"

	"github.com/costa92/go-web/internal/middleware"
	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/util"
)

type MenuListRequest struct {
	model.PageRequest
	MenuStatus int `json:"menu_status" query:"menu_status"  form:"menu_status"`
}

func (m *MenuController) List(ctx *gin.Context) {
	var req MenuListRequest
	ret := &model.MenuList{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		util.WriteErrResponse(ctx, code.ErrBind, err, nil)
		return
	}
	tx := m.MysqlStorage.
		Model(&model.Menu{}).Debug()
	if req.MenuStatus > 0 {
		tx = tx.Where("menu_status = ?", req.MenuStatus)
	}
	// 处理查询数据
	if err := tx.Order("menu_sort asc").
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
