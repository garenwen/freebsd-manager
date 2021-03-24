package client

import (
	"fmt"

	v1 "github.com/garenwen/freebsd-manager/pkg/apis/storage/v1"
)

type Config struct {
	Url    string
	APIVer string
}

type response struct {
	*HttpRequest
	res v1.BaseResult
}

func NewResponse(r *HttpRequest) *response {
	return &response{
		HttpRequest: r,
	}
}

func (r *response) IntoBaseRes(data interface{}) error {
	if nil != data {
		r.res = v1.BaseResult{
			Data: data,
		}
	}
	if err := r.IntoJson(&r.res); nil != err {
		return err
	}
	if r.res.Status != v1.StatusSuccess {
		return fmt.Errorf("converter response error: errorType：%s msg: %s", r.res.ApiError.Typ, r.res.ApiError.Msg)
	}

	return nil

}
