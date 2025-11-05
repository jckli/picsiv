package utils

import (
	"github.com/valyala/fasthttp"
)

var (
	client = &fasthttp.Client{}
)

func getRequest(url string) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetMethod("GET")
	req.SetRequestURI(url)

	resp := fasthttp.AcquireResponse()

	if err := client.Do(req, resp); err != nil {
		return nil, err
	}

	if resp.StatusCode() != fasthttp.StatusOK {
	}

	return resp.Body(), nil
}
