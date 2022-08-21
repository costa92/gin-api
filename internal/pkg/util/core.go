package util

import (
	"fmt"
	"net/http"

	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"
)

type ErrResponse struct {
	Code int `json:"code"`

	Message string `json:"message"`

	Reference string `json:"reference"`
}

type SuccessResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

func WriteResponse(c *gin.Context, err error, data interface{}) {
	if err != nil {
		fmt.Printf("%#+v", err)
		coder := errors.ParseCoder(err)
		c.JSON(coder.HTTPStatus(), ErrResponse{
			Code:      coder.Code(),
			Message:   coder.String(),
			Reference: coder.Reference(),
		})
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{
		Code:    http.StatusOK,
		Message: "Success",
		Result:  data,
	})
}
