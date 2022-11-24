package enterprise

import (
	"time"

	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"

	"github.com/costa92/go-web/internal/middleware"
	"github.com/costa92/go-web/internal/validation"
	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/util"
)

func (c *EnterpriseController) Create(ctx *gin.Context) {
	var req CreateRequest
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
	enterprise := &model.Enterprise{}
	enterpriseModel := model.NewEnterpriseModel(ctx, c.MysqlStorage)
	c.saveParams(enterprise, &req)
	currUserId := middleware.GetAuthUserId(ctx)
	enterprise.CreatedBy = currUserId
	enterprise.UpdateBy = currUserId
	if err := enterpriseModel.Save(enterprise); err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()),
			"企业添加类型数据错误")
		return
	}
	if len(req.Contacts) > 0 {
		currTime := time.Now().Unix()
		contactModel := model.NewEnterpriseContactModel(ctx, c.MysqlStorage)
		createContacts := make([]*model.EnterpriseContact, 0)
		for _, contact := range req.Contacts {
			createContacts = append(createContacts, &model.EnterpriseContact{
				Name:         contact.Name,
				Mobile:       contact.Mobile,
				Position:     contact.Position,
				EnterpriseID: enterprise.Id,
				CreatedBy:    currUserId,
				UpdatedBy:    currUserId,
				UpdatedAt:    currTime,
				CreatedAt:    currTime,
				Status:       model.ContactStatusNormal,
			})
		}

		if len(createContacts) > 0 {
			if err := contactModel.DB.Create(createContacts).Error; err != nil {
				util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "修改企业联系人错误")
				return
			}
		}
	}

	util.WriteSuccessResponse(ctx, "success")
}
