package enterprise

import (
	"fmt"

	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"

	"github.com/costa92/go-web/internal/validation"
	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/util"
)

func (c *EnterpriseController) Update(ctx *gin.Context) {
	var req UpdateRequest
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
	enterprise, err := enterpriseModel.FirstById(req.Id)
	if err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "查询企业数据错误")
		return
	}

	c.saveParams(enterprise, &req.CreateRequest)
	if err := enterpriseModel.Save(enterprise); err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()),
			"企业修改类型数据错误")
		return
	}
	contactModel := model.NewEnterpriseContactModel(ctx, c.MysqlStorage)

	if len(req.Contacts) > 0 {
		// 查询历史联系人
		contacts, err := contactModel.FindByEnterpriseId(int(enterprise.Id))
		if err != nil {
			util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "查询企业联系人错误")
			return
		}

		updateIds := map[int32]*ContactItems{}
		for _, contact := range req.Contacts {
			if contact.Id > 0 {
				updateIds[contact.Id] = contact
			}
		}

		updateContacts := make([]*model.EnterpriseContact, 0)
		contactItems := make(map[int32]*model.EnterpriseContact, 0)
		for _, contact := range contacts {
			contactItems[contact.ID] = contact
			if ok := updateIds[contact.ID]; ok == nil {
				item := contactItems[contact.ID]
				item.Status = model.ContactStatusFail
				fmt.Println(contact.ID)
				updateContacts = append(updateContacts, item)
			}
		}

		createContacts := make([]*model.EnterpriseContact, 0)
		for _, contact := range req.Contacts {
			if contact.Id > 0 {
				if ok := contactItems[contact.Id]; ok != nil {
					// 修改的
					item := contactItems[contact.Id]
					item.Name = contact.Name
					item.Mobile = contact.Mobile
					item.Position = contact.Position
					updateContacts = append(updateContacts, item)
				}
			} else {
				// 新增加的
				createItem := &model.EnterpriseContact{
					Name:         contact.Name,
					Mobile:       contact.Mobile,
					Position:     contact.Position,
					EnterpriseID: enterprise.Id,
					Status:       model.ContactStatusNormal,
				}
				createContacts = append(createContacts, createItem)
			}
		}
		if len(updateContacts) > 0 {
			if err := contactModel.DB.Save(updateContacts).Error; err != nil {
				util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "修改企业联系人错误")
				return
			}
		}
		if len(createContacts) > 0 {
			if err := contactModel.DB.Create(createContacts).Error; err != nil {
				util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "修改企业联系人错误")
				return
			}
		}
	} else { // 如果没有联系人
		if err := contactModel.DB.Where("enterprise_id = ?", enterprise.Id).
			UpdateColumns(model.EnterpriseContact{Status: model.ContactStatusFail}).Error; err != nil {
			util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), "修改企业联系人错误")
			return
		}
	}

	util.WriteSuccessResponse(ctx, "success")
}
