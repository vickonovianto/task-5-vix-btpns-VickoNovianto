package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	response struct {
		Status  bool        `json:"status"`
		Message string      `json:"message"`
		Errors  []string    `json:"errors"`
		Data    interface{} `json:"data"`
	}
)

func SendResponseSuccess(c *gin.Context, data interface{}) {
	message := "Success"
	res := response{
		Status:  true,
		Message: message,
		Errors:  nil,
		Data:    data,
	}
	c.JSON(http.StatusOK, res)
	c.Abort()
}

func SendResponseError(c *gin.Context, code int, err error) {
	message := "Failed"
	res := response{
		Status:  false,
		Message: message,
		Errors:  []string{err.Error()},
		Data:    nil,
	}
	c.JSON(code, res)
	c.Abort()
}
