package v1

import (
	"testing"

	v1 "github.com/garenwen/freebsd-manager/pkg/apis/storage/v1"
	"github.com/garenwen/freebsd-manager/pkg/client"
)

const (
	urls       = "http://127.0.0.1:8870"
	apiVersion = "v1"
)

func TestClient_CreateVolume(t *testing.T) {
	crResp, err := NewClientOrDie(client.Config{
		Url:    urls,
		APIVer: apiVersion,
	}).CreateVolume(&v1.CreateVolumeRequest{
		Name: "test/appv1",
		CapacityRange: &v1.CapacityRange{
			RequiredBytes: 1024,
			LimitBytes:    1024,
		},
	})
	if err != nil {
		t.Errorf("err: %v \n", err)
		return
	}

	t.Log("success")
	t.Logf("crResp %#v \n", crResp)

}
