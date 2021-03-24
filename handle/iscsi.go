package handle

import (
	"net/http"

	v1 "github.com/garenwen/freebsd-manager/pkg/apis/storage/v1"
	"github.com/gin-gonic/gin"
)

type IscsiHandler struct {
}

func NewIscsiHandler() *IscsiHandler {
	return &IscsiHandler{}
}

func (ih *IscsiHandler) HandleCreateIscsi(c *gin.Context) {
	c.JSON(http.StatusOK, ih.createIscsi(c))
}

func (ih *IscsiHandler) createIscsi(c *gin.Context) v1.BaseResult {

	return v1.BaseResult{Status: v1.StatusSuccess, Data: nil, ApiError: nil}
}
