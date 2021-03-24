package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// Class for HTTP request
type HttpRequest struct {
	// Error if any. Will be checked at first before any method execution.
	err error

	// URL info to where to send HTTP request
	url *url.URL
	// HTTP method
	method string
	// HTTP request headers
	headers http.Header
	// HTTP query parameters
	query url.Values
	// HTTP request body
	reqBody io.Reader
	// HTTP response
	resp *http.Response
	// Debug flag
	debug bool
}

// New HttpRequest object with specified URL
func NewHttpReq(u url.URL) *HttpRequest {
	return &HttpRequest{
		url:     &u,
		headers: make(http.Header),
	}
}

// Set path of URL
func (r *HttpRequest) BasicAuth(username, password string) *HttpRequest {
	if nil == r.err {
		r.headers.Set("Authorization", authorizationHeader("measure", "measure"))
	}

	return r
}
func authorizationHeader(user, password string) string {
	base := user + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(base))
}

// Set path of URL
func (r *HttpRequest) Path(path string) *HttpRequest {
	if nil == r.err {
		r.url.Path = path
	}

	return r
}

// Append sub path of URL
func (r *HttpRequest) SubPath(subPath string) *HttpRequest {
	if nil == r.err {
		r.url.Path = path.Join(r.url.Path, subPath)
	}

	return r
}

func (r *HttpRequest) RawSubPath(subPath string) *HttpRequest {
	if nil == r.err {
		r.url.RawPath = path.Join(r.url.Path, subPath)

		subPath, r.err = url.PathUnescape(subPath)
		if nil == r.err {
			r.url.Path = path.Join(r.url.Path, subPath)
		}
	}

	return r
}

// Set HTTP method
func (r *HttpRequest) Method(m string) *HttpRequest {
	if nil == r.err {
		r.method = m
	}

	return r
}

// Set specified header
func (r *HttpRequest) SetHeader(k, v string) *HttpRequest {
	if nil == r.err {
		if nil == r.headers {
			r.headers = make(http.Header)
		}
		r.headers.Set(k, v)
	}

	return r
}

// Add value to specified header
func (r *HttpRequest) AddHeader(k, v string) *HttpRequest {
	if nil == r.err {
		if nil == r.headers {
			r.headers = make(http.Header)
		}
		r.headers.Add(k, v)
	}

	return r
}

// Set query of URL
func (r *HttpRequest) Query(params url.Values) *HttpRequest {
	if nil == r.err {
		r.url.RawQuery = params.Encode()
	}

	return r
}

// Set request body
func (r *HttpRequest) Body(body io.Reader) *HttpRequest {
	if nil == r.err {
		r.reqBody = body
	}

	return r
}

// Set form body with specified kv pairs
func (r *HttpRequest) FormBody(vals url.Values) *HttpRequest {
	if nil == r.err {
		r.headers.Set("Content-Type", "application/x-www-form-urlencoded")
		body := vals.Encode()
		if r.debug {
			log.Println("form-urlencoded request body:", body)
		}
		r.reqBody = strings.NewReader(body)
	}

	return r
}

// Set JSON encoded body with specified object
func (r *HttpRequest) JsonBody(obj interface{}) *HttpRequest {
	if nil != r.err {
		return r
	}

	body := bytes.NewBuffer(nil)
	if err := json.NewEncoder(body).Encode(obj); nil != err {
		r.err = fmt.Errorf("Encode request body into JSON error: %v \n", err)
		return r
	}

	r.reqBody = body
	if r.debug {
		log.Println("Json encoded request body:", body.String())
	}
	r.headers.Set("Content-Type", "application/json")
	return r
}

func (r *HttpRequest) Debug() *HttpRequest {
	r.debug = true
	return r
}

// Do send request and return proxy.Response object and error.
// The proxy.Response.Body should be closed by caller
func (r *HttpRequest) DoRaw() (*http.Response, error) {
	if nil != r.err {
		return nil, r.err
	}
	u := r.url.String()

	if r.debug {
		log.Println("request url:", u)
		log.Println("request headers:", r.headers)
		log.Println("request method:", r.method)
	}

	req, err := http.NewRequest(r.method, u, r.reqBody)
	if nil != err {
		return nil, fmt.Errorf("Init request error: %v \n", err)
	}

	req.Header = r.headers

	resp, err := http.DefaultClient.Do(req)
	if nil != err {
		return nil, fmt.Errorf("Request remote error: %v \n", err)
	}

	return resp, nil
}

// Do send request.
// One of the following methods should be called after this method:
//   * Call Error() to get error
//   * IntoJson to extract JSON encoded response body
func (r *HttpRequest) Do() *HttpRequest {

	r.resp, r.err = r.DoRaw()
	return r
}

// Get current error
func (r *HttpRequest) Error() error {
	if nil != r.resp {
		r.resp.Body.Close()
	}
	return r.err
}

// Decode response body into specified object.
// Do() should be called before
func (r *HttpRequest) IntoJson(expected interface{}) error {
	if nil != r.err || nil == r.resp {
		return r.err
	}

	if nil != r.resp {
		// Close body
		defer r.resp.Body.Close()
		// Read data
		data, _ := ioutil.ReadAll(r.resp.Body)

		if nil == expected {
			// Return error if no expected object and status code >= 400
			if r.resp.StatusCode >= http.StatusBadRequest {
				return fmt.Errorf("Service response error: (%d)%s", r.resp.StatusCode, string(data))
			}
			return nil
		}

		if err := json.Unmarshal(data, expected); nil != err {
			if r.resp.StatusCode >= http.StatusBadRequest {
				return fmt.Errorf("Service response error: (%d)%s", r.resp.StatusCode, string(data))
			}
			return fmt.Errorf("Decode response error: %v", err)
		}
	}

	return nil
}
