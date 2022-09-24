package menus

import (
	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"

	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/util"
)

type ListRequest struct{}

func (m *MenuController) List(ctx *gin.Context) {
	var req ListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrBind, err.Error()), nil)
		return
	}
	util.WriteResponse(ctx, nil, "ret")
}

type MenuCreateRequest struct {
	ParentID      uint   `json:"parent_id"`
	MenuName      string `json:"menu_name"`
	MenuType      uint   `json:"menu_type"`
	MenuSort      uint   `json:"menu_sort"`
	MenuStatus    int    `json:"menu_status"`
	HideMenu      int    `json:"hide_menu"`
	ExternalLink  int    `json:"external_link"`
	Permission    string `json:"permission"`
	ComponentName string `json:"component_name"`
	Component     string `json:"component"`
	KeepAlive     uint   `json:"keep_alive"`
	Icon          string `json:"icon"`
	Path          string `json:"path"`
}

type MenuUpdateRequest struct {
	ID int `json:"id" binding:"required"`
	MenuCreateRequest
}

func (m *MenuController) SaveParams(menu *model.Menu, req *MenuCreateRequest) *model.Menu {
	menu.ParentID = req.ParentID
	menu.MenuName = req.MenuName
	menu.MenuType = req.MenuType
	menu.MenuSort = req.MenuSort
	menu.MenuStatus = req.MenuStatus
	menu.HideMenu = req.HideMenu
	menu.ExternalLink = req.ExternalLink
	menu.Permission = req.Permission
	menu.ComponentName = req.ComponentName
	menu.Component = req.Component
	menu.KeepAlive = req.KeepAlive
	menu.Icon = req.Icon
	menu.Path = req.Path
	return menu
}
