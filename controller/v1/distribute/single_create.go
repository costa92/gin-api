package distribute

import (
	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/costa92/go-web/internal/validation"
	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/util"
)

type SingleCreateRequest struct {
	UserId       int32 `json:"user_id" form:"user_id" validate:"required"`
	EnterpriseId uint  `json:"enterprise_id" form:"enterprise_id" validate:"required"`
}

func (c *DistributeController) SingleCreate(ctx *gin.Context) {
	var req SingleCreateRequest
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
	enterpriseModel := model.NewEnterpriseModel(ctx, c.MysqlStorage)
	tx := enterpriseModel.DB
	enterpriseInfo, err := enterpriseModel.FirstById(int(req.EnterpriseId))
	if err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "查询数据错误")
		return
	}
	if enterpriseInfo.Status != 1 {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, "该数据已经处理"), "数据已经分发")
		return
	}
	distributeModel := model.NewDistributeModel(ctx, c.MysqlStorage)
	distributeInfo, err := distributeModel.FirstByEnterpriseID(req.EnterpriseId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "该数据已经存在")
		return
	}
	if distributeInfo != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrValidation, "该数据已经分发"), "该数据已经分发")
		return
	}
	distribute := &model.Distribute{
		UserID:       req.UserId,
		Status:       1,
		EnterpriseID: req.EnterpriseId,
	}
	if err := distributeModel.Save(distribute); err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "保存数据错误")
		return
	}
	enterpriseInfo.Status = 2
	if err := enterpriseModel.Save(enterpriseInfo); err != nil {
		tx.Rollback()
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "保存数据错误")
		return
	}
	tx.Commit()
	util.WriteSuccessResponse(ctx, "success")
}
