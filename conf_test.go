//Copyright (c) 2017 Phil
package apollo

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type TestConfigSuite struct {
	suite.Suite
}

func (ts *TestConfigSuite) TestNewConf() {
	var tcs = []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "fakename",
			wantErr: true,
		},
		{
			name:    "./LICENSE",
			wantErr: true,
		},
		{
			name:    "./testdata/" + defaultConfName,
			wantErr: false,
		},
	}

	for _, tc := range tcs {
		_, err := NewConf(tc.name)
		if tc.wantErr {
			ts.Error(err)
		} else {
			ts.NoError(err)
		}
	}
}

func TestRunConfigSuite(t *testing.T) {
	suite.Run(t, new(TestConfigSuite))
}
