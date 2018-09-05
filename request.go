//Copyright (c) 2017 Phil

package apollo

import (
	"io"
	"io/ioutil"
	"net/http"
)

// this is a static check
var _ requester = (*httpRequester)(nil)

type requester interface {
	request(url string) ([]byte, error)
}

type httpRequester struct {
	client *http.Client
}

func newHTTPRequester(client *http.Client) requester {
	return &httpRequester{
		client: client,
	}
}

func (r *httpRequester) request(url string) ([]byte, error) {
	resp, err := r.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return ioutil.ReadAll(resp.Body)
	}

	// Discard all body if status code is not 200
	io.Copy(ioutil.Discard, resp.Body)
	return nil, nil
}
