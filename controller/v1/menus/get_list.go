package menus

import (
	"github.com/costa92/errors"
	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/util"
	"github.com/gin-gonic/gin"
)

type RouteMeta struct {
	Title           string `json:"title"`
	Icon            string `json:"icon"`
	OrderNo         uint   `json:"orderNo"`
	IgnoreKeepAlive bool   `json:"ignoreKeepAlive"` // true 即可关闭缓存
	Affix           bool   `json:"affix"`
	HideMenu        bool   `json:"hideMenu"`
}

type MenuTreeItem struct {
	ID        uint           `json:"id"`
	Path      string         `json:"path"`
	Name      string         `json:"name"`
	Component string         `json:"component"`
	Meta      RouteMeta      `json:"meta"`
	Children  []MenuTreeItem `json:"children"`
}
type GetMenuListRequest struct {
}

func (m *MenuController) GetMenuList(ctx *gin.Context) {
	var menus []*model.Menu
	if err := m.MysqlStorage.Model(&model.Menu{}).Where("menu_type in ?", []int{1, 2}).Order("id desc").Find(&menus).Error; err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "查询数据库错误")
		return
	}
	trees := toMenuTree(menus, 0)
	util.WriteSuccessResponse(ctx, trees)
	return
}

func toMenuTree(menus []*model.Menu, parentID uint) []MenuTreeItem {
	tree := make([]MenuTreeItem, 0)
	for _, menu := range menus {
		if parentID != menu.ParentID {
			continue
		}
		component := menu.Component
		if menu.MenuType == 1 {
			component = "LAYOUT"
		}
		componentName := menu.MenuName
		if menu.ComponentName != "" {
			componentName = menu.ComponentName
		}
		item := MenuTreeItem{
			ID:        menu.ID,
			Path:      menu.Path,
			Name:      componentName,
			Component: component,
			Meta: RouteMeta{
				Title:           menu.MenuName,
				Icon:            menu.Icon,
				OrderNo:         menu.MenuSort,
				IgnoreKeepAlive: menu.KeepAlive == 0,
				Affix:           false,
				HideMenu:        menu.HideMenu == 2,
			},
			Children: []MenuTreeItem{},
		}
		item.Children = toMenuTree(menus, menu.ID)
		tree = append(tree, item)
	}
	return tree
}
