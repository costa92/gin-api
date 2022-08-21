package index

import (
	"fmt"

	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"github.com/costa92/go-web/internal/pkg/code"
	"github.com/costa92/go-web/internal/pkg/meta"
	"github.com/costa92/go-web/internal/pkg/util"
	"github.com/costa92/go-web/internal/pkg/util/gormutil"
	"github.com/costa92/go-web/model"
)

type Index struct {
	MysqlStorage *gorm.DB
}

func NewIndex(db *gorm.DB) *Index {
	return &Index{MysqlStorage: db}
}

func (api *Index) Index(ctx *gin.Context) {
	var r meta.ListOptions
	ret := &model.AdminList{}
	if err := ctx.ShouldBindQuery(&r); err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrBind, err.Error()), nil)
		return
	}
	// 分页处理
	ol := gormutil.Unpointer(r.Offset, r.Limit)
	// 处理查询数据
	if err := api.MysqlStorage.
		Model(&model.Admin{}).Scopes(gormutil.Paginate(ol.Offset, ol.Limit)).
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

type SignUpParams struct {
	Age uint8 `form:"age" binding:"gte=1,lte=130" json:"age"`
}

func (api *Index) Create(ctx *gin.Context) {
	var sp SignUpParams
	if err := ctx.ShouldBind(&sp); err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrValidation, err.Error()), "参数错误")
		return
	}
	fmt.Println(&sp)
	util.WriteResponse(ctx, nil, "data")
}

type UpSignUpParams struct {
	SignUpParams
	Id int `form:"id" binding:"gte=1,lte=130" json:"id"`
}

func (api *Index) Update(ctx *gin.Context) {
	var sp SignUpParams
	if err := ctx.ShouldBind(&sp); err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrValidation, err.Error()), "参数错误")
		return
	}
}

type RequestDetail struct {
	ID int `query:"id" binding:"gte=1" form:"id"`
}

func (api *Index) Detail(ctx *gin.Context) {
	var req RequestDetail
	if err := ctx.ShouldBind(&req); err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrValidation, err.Error()), "参数错误")
		return
	}
	util.WriteResponse(ctx, nil, "ddd")
}

func (api *Index) Del(ctx *gin.Context) {
	log.Info().Msg("del")
}
