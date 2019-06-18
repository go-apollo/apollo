//Copyright (c) 2017 Phil

package apollo

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/suite"
)

type CommonTestSuite struct {
	suite.Suite
}

func (s *CommonTestSuite) TestLocalIp() {
	ip := getLocalIP()
	s.NotEmpty(ip)
}

func (s *CommonTestSuite) TestNotificationURL() {
	target := notificationURL(
		&Conf{
			IP:      "127.0.0.1:8080",
			AppID:   "SampleApp",
			Cluster: "default",
		}, "")
	_, err := url.Parse(target)
	s.NoError(err)
}

func (s *CommonTestSuite) TestConfigURL() {
	target := configURL(
		&Conf{
			IP:      "127.0.0.1:8080",
			AppID:   "SampleApp",
			Cluster: "default",
		}, "application", "")
	_, err := url.Parse(target)
	s.NoError(err)
}

func TestRunCommonSuite(t *testing.T) {
	cs := new(CommonTestSuite)
	suite.Run(t, cs)
}
