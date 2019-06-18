//Copyright (c) 2017 Phil

package apollo

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TestRequestSuite struct {
	suite.Suite
}

func (ts *TestRequestSuite) TestRequest() {
	request := newHTTPRequester(&http.Client{})

	serv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("test"))
	}))

	bts, err := request.request(serv.URL)
	ts.NoError(err)

	ts.Equal(bts, []byte("test"))

	serv = httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusInternalServerError)
	}))
	bts, err = request.request(serv.URL)
	ts.NoError(err)

	ts.Empty(bts)

	serv = httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusInternalServerError)
	}))
	serv.Close()
	_, err = request.request(serv.URL)
	ts.Error(err)
}

func TestRunRequestSuite(t *testing.T) {
	ts := new(TestRequestSuite)
	suite.Run(t, ts)
}
