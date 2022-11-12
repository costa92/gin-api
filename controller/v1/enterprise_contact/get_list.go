package enterprise_contact

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

type GetListRequest struct {
	EnterpriseId int `json:"enterprise_id,omitempty" form:"enterprise_id" query:"enterprise_id"`
}

type GetListItem struct {
	*model.EnterpriseContact
	CreatedTime string `json:"created_time"`
}

func (c *EnterpriseContactController) GetListByEnterpriseId(ctx *gin.Context) {
	var req GetListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrBind, err.Error()), nil)
		return
	}
	validate := validation.DefaultValidator{}
	err := validate.ValidateStruct(req)
	if err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrValidation, err.Error()), nil)
		return
	}
	contactModel := model.NewEnterpriseContactModel(ctx, c.MysqlStorage)
	contacts, err := contactModel.FindByEnterpriseId(req.EnterpriseId)
	if err != nil {
		util.WriteResponse(ctx, errors.WithCode(code.ErrDatabase, err.Error()), nil)
	}
	items := make([]*GetListItem, 0)
	if len(contacts) > 0 {
		for _, contact := range contacts {
			items = append(items, &GetListItem{
				EnterpriseContact: contact,
				CreatedTime:       time.Unix(contact.CreatedAt, 0).Format(middleware.TimeFieldFormat),
			})
		}
	}
	util.WriteSuccessResponse(ctx, items)
}
