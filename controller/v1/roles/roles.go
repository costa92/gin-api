package roles

import (
	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/meta"
	"github.com/costa92/go-web/pkg/util"
	"github.com/costa92/go-web/pkg/util/gormutil"
)

type RoleController struct {
	MysqlStorage *gorm.DB
}

func NewRoleController(db *gorm.DB) *RoleController {
	return &RoleController{MysqlStorage: db}
}

func (api *RoleController) Index(ctx *gin.Context) {
	var r meta.ListOptions
	ret := &model.RoleList{}
	if err := ctx.ShouldBindQuery(&r); err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrBind, err.Error()), nil)
		return
	}
	// 分页处理
	ol := gormutil.Unpointer(r.Offset, r.Limit)
	// 处理查询数据
	if err := api.MysqlStorage.
		Model(&model.Role{}).Scopes(gormutil.Paginate(ol.Offset, ol.Limit)).
		Order("id desc").
		Find(&ret.Items).
		Offset(-1).
		Limit(-1).
		Count(&ret.TotalCount).
		Error; err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "查询数据库错误")
		return
	}
	util.WriteResponse(ctx, nil, ret)
}

type CreateRequest struct {
	Name   string `form:"Name" binding:"required" json:"name"`
	Remark string `form:"remark" binding:"required" json:"remark"`
}

func (api *RoleController) Create(ctx *gin.Context) {
	var rep CreateRequest
	if err := ctx.ShouldBind(&rep); err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrValidation, err.Error()), "参数错误")
		return
	}
	role := &model.Role{
		Name:   rep.Name,
		Remark: rep.Remark,
	}
	if err := api.MysqlStorage.Model(&model.Role{}).Create(role).Error; err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "添加角色错误")
		return
	}
	util.WriteResponse(ctx, nil, "success")
}

type UpdateRequest struct {
	CreateRequest
	Id int `form:"id" binding:"gte=1" json:"id"`
}

func (api *RoleController) Update(ctx *gin.Context) {
	var req UpdateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrValidation, err.Error()), "参数错误")
		return
	}
	var role model.Role
	tx := api.MysqlStorage.Model(&model.Role{})
	if err := tx.Where("id = ?", req.Id).First(&role).Error; err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "查询数据错误")
		return
	}
	role.Name = req.Name
	role.Remark = req.Remark
	if err := tx.Save(&role).Error; err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "保存数据错误")
		return
	}
	util.WriteResponse(ctx, nil, "success")
}

type RequestDetail struct {
	Id int `query:"id" binding:"gte=1" form:"id"`
}

func (api *RoleController) Detail(ctx *gin.Context) {
	var req RequestDetail
	if err := ctx.ShouldBind(&req); err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrValidation, err.Error()), "参数错误")
		return
	}
	var role model.Role
	tx := api.MysqlStorage.Model(&model.Role{})
	if err := tx.Where("id = ?", req.Id).First(&role).Error; err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "查询数据错误")
		return
	}
	util.WriteResponse(ctx, nil, role)
}

func (api *RoleController) Del(ctx *gin.Context) {
	log.Info().Msg("del")
}

type UpdateStateRequest struct {
	Id     int `json:"id" binding:"gte=1" form:"id"`
	Status int `json:"status" binding:"required,oneof=1 2" form:"status"`
}

func (api *RoleController) UpdateState(ctx *gin.Context) {
	var req UpdateStateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrValidation, err.Error()), "参数错误")
		return
	}
	var role model.Role
	tx := api.MysqlStorage.Model(&model.Role{})
	if err := tx.Where("id = ?", req.Id).First(&role).Error; err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "查询数据错误")
		return
	}
	role.Status = req.Status
	if err := tx.Save(&role).Error; err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "修改角色状态错误")
		return
	}
	util.WriteResponse(ctx, nil, "success")
}
