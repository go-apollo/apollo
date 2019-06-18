//Copyright (c) 2017 Phil

package apollo

import (
	"log"
	"os"
	"time"

	// "time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/apollo.v0/internal/mockserver"
)

type StartWithConfTestSuite struct {
	suite.Suite
	changeEvent <-chan *ChangeEvent
	// updates = WatchUpdate()
	// apollo *Client
	// VariableThatShouldStartAtFive int
}

func (s *StartWithConfTestSuite) SetupTest() {

	if err := StartWithConfFile("./testdata/app.yml"); err != nil {
		log.Fatalf("Start with app.yml should return nil, got :%v", err)
	}
	// defer Stop()

	if err := defaultClient.loadLocal(defaultDumpFile); err != nil {
		log.Fatalf("loadLocal should return nil, got: %v", err)
	}

	s.changeEvent = WatchUpdate()

	// s.apollo = defaultClient
}
func (s *StartWithConfTestSuite) TearDownTest() {
	Stop()
	defer os.Remove(defaultDumpFile)
}
func (s *StartWithConfTestSuite) LoadLocal() {
	err := defaultClient.loadLocal(defaultDumpFile)
	assert.NoError(s.T(), err)
}
func (s *StartWithConfTestSuite) BeforeGetStringValueWithNameSpace() {
	mockserver.Set("application", "key", "value")
	s.wait()
}
func (s *StartWithConfTestSuite) GetStringValueWithNameSpace() {
	val := GetStringValueWithNameSpace("application", "key", "defaultValue")
	assert.Equal(s.T(), "value", val)
}
func (s *StartWithConfTestSuite) BeforeGetStringValue() {
	mockserver.Set("application", "key", "newvalue")
	s.wait()
}
func (s *StartWithConfTestSuite) GetStringValue() {
	val := GetIntValue("intkey", 0)
	assert.Equal(s.T(), 1, val)
}

func (s *StartWithConfTestSuite) BeforeGetIntValue() {
	mockserver.Set("application", "key", "newvalue")
	s.wait()
}
func (s *StartWithConfTestSuite) GetIntValue() {
	val := GetStringValue("key", "defaultValue")
	assert.Equal(s.T(), "newvalue", val)
}

func (s *StartWithConfTestSuite) BeforeGetNameSpaceContent() {

	mockserver.Set("client.json", "content", `{"name":"apollo"}`)
	s.wait()
}
func (s *StartWithConfTestSuite) GetNameSpaceContent() {

	val := GetNameSpaceContent("client.json", "{}")
	assert.Equal(s.T(), `{"name":"apollo"}`, val)
}

func (s *StartWithConfTestSuite) wait() {
	select {
	case <-s.changeEvent:
	case <-time.After(time.Second * 30):
	}
}
