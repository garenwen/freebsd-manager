package handle

import (
	"net/http"

	v1 "github.com/garenwen/freebsd-manager/pkg/apis/storage/v1"

	"github.com/garenwen/freebsd-manager/pkg/zfs"

	"github.com/gin-gonic/gin"
)

type ZfsHandler struct {
}

func NewZfsHandler() *ZfsHandler {

	return &ZfsHandler{}
}
func (zfsHandler *ZfsHandler) HandleCreateVolume(c *gin.Context) {

	result := zfsHandler.createVolume(c)
	c.JSON(http.StatusOK, result)
}

func (zfsHandler *ZfsHandler) createVolume(c *gin.Context) v1.BaseResult {

	cvr := v1.CreateVolumeRequest{}
	if err := c.ShouldBindJSON(&cvr); err != nil {
		return v1.BaseResult{Status: v1.StatusError, ApiError: &v1.ApiError{Typ: v1.ErrorBadData, Msg: err.Error()}}
	}
	v, err := zfs.CreateVolume(cvr.Name, uint64(cvr.CapacityRange.RequiredBytes), cvr.Parameters)
	if err != nil {
		return v1.BaseResult{Status: v1.StatusError, ApiError: &v1.ApiError{Typ: v1.ErrorBadData, Msg: err.Error()}}
	}

	return v1.BaseResult{Status: v1.StatusSuccess, Data: v, ApiError: nil}
}
