package distribute

import (
	"sync"

	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/costa92/go-web/internal/validation"
	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/logger"
	"github.com/costa92/go-web/pkg/util"
)

type AreaCreateRequest struct {
	ParentCode []int `json:"parentCode,omitempty" form:"parentCode" validate:"required"`
	UserId     int32 `json:"user_id" form:"user_id" validate:"required"`
	Confirm    bool  `json:"confirm" form:"confirm"`
}

func (c *DistributeController) AreaCreate(ctx *gin.Context) {
	var req AreaCreateRequest
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
	var provinceId, cityId, countyId int
	if len(req.ParentCode) > 0 {
		provinceId = req.ParentCode[0]
		cityId = req.ParentCode[1]
		countyId = req.ParentCode[2]
	}
	enterpriseModel := model.NewEnterpriseModel(ctx, c.MysqlStorage)
	if req.Confirm == false {
		total, err := enterpriseModel.CountByArea(provinceId, cityId, countyId, 1)
		if err != nil {
			util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), nil)
			return
		}
		util.WriteSuccessResponse(ctx, total)
		return
	} else {
		enterprises, err := enterpriseModel.FindByArea(provinceId, cityId, countyId, 1)
		if err != nil {
			util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), nil)
			return
		}
		distribute := model.Distribute{
			UserID: req.UserId,
			Status: 1,
		}
		distributeModel := model.NewDistributeModel(ctx, c.MysqlStorage)
		var wg sync.WaitGroup
		var enterpriseIds []uint
		for _, enterprise := range enterprises {
			wg.Add(1)
			go func(enterprise *model.Enterprise) {
				defer func() {
					if err := recover(); err != nil {
						logger.Errorw("AreaCreate", "distribute", distribute, "err", err)
					}
					wg.Done()
				}()

				distributeInfo, err := distributeModel.FirstByEnterpriseID(enterprise.Id)
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					logger.Errorw("distributeModel.FirstByEnterpriseID distributeInfo", "distributeInfo", distributeInfo)
				} else {
					distribute.EnterpriseID = enterprise.Id
					logger.Infow("AreaCreate distribute", "distribute", distribute)
					err := distributeModel.Save(&distribute)
					if err != nil {
						logger.Errorw("AreaCreate distribute save", "distribute", distribute, "err", err)
					} else {
						// 修改企业状态
						enterprise.Status = 2
						if err := enterpriseModel.DB.Where("id =?", enterprise.Id).Debug().Updates(enterprise).Error; err != nil {
							logger.Errorw("AreaCreate enterpriseModel save", "distribute", distribute, "err", err, "enterprise", enterprise)
						}
					}
				}
			}(enterprise)
		}
		wg.Wait()
		util.WriteSuccessResponse(ctx, enterpriseIds)
		return
	}
}
