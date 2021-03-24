package v1

import (
	"fmt"
	"net/http"
	"net/url"

	v1 "github.com/garenwen/freebsd-manager/pkg/apis/storage/v1"
	"github.com/garenwen/freebsd-manager/pkg/client"
)

const (
	createVolume = "/create_volume"
)

type Client struct {
	url    *url.URL `env:"URL"`
	apiVer string
}

func NewClientOrDie(conf client.Config) *Client {
	c, err := NewClient(conf)
	if nil != err {
		panic(fmt.Errorf("new client error: %v", err))
	}

	return c
}

func NewClient(conf client.Config) (*Client, error) {
	u, err := url.Parse(conf.Url)
	if nil != err {
		return nil, fmt.Errorf("parse external url error: %v", err)
	}

	if "" == conf.APIVer {
		conf.APIVer = "v1"
	}

	return &Client{
		url:    u,
		apiVer: conf.APIVer,
	}, nil
}

func (c *Client) newRequest() *client.HttpRequest {
	return client.NewHttpReq(*c.url).Path("apis/storage/" + c.apiVer)
}

// Query namespace based on username or userid.
func (c *Client) CreateVolume(req *v1.CreateVolumeRequest) (*v1.CreateVolumeResponse, error) {

	var cvResp *v1.CreateVolumeResponse
	hr := c.newRequest().Debug().
		Method(http.MethodPost).
		SubPath(createVolume).
		JsonBody(req).
		Do()

	if err := client.NewResponse(hr).
		IntoBaseRes(&cvResp); err != nil {
		return nil, err
	}

	return cvResp, nil
}
