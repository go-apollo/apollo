//Copyright (c) 2017 Phil

package apollo

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"gopkg.in/apollo.v0/internal/mockserver"
)

type StartWithConfTestSuite struct {
	suite.Suite
	changeEvent <-chan *ChangeEvent
}

func (s *StartWithConfTestSuite) SetupSuite() {
	startTestApollo()
	s.changeEvent = WatchUpdate()
}
func startTestApollo() {
	if err := StartWithConfFile("./testdata/app.yml"); err != nil {
		log.Fatalf("Start with app.yml should return nil, got :%v", err)
	}
	// defer Stop()

	if err := defaultClient.loadLocal(defaultDumpFile); err != nil {
		log.Fatalf("loadLocal should return nil, got: %v", err)
	}
}
func (s *StartWithConfTestSuite) BeforeTest(suiteName, testName string) {
	// log.Println("suiteName: " + suiteName)
	// log.Println("testName: " + testName)
}
func (s *StartWithConfTestSuite) TearDownSuite() {
	Stop()
	os.Remove(defaultDumpFile)
}
func (s *StartWithConfTestSuite) TestLoadLocal() {
	err := defaultClient.loadLocal(defaultDumpFile)
	s.NoError(err)
}

func (s *StartWithConfTestSuite) TestGetStringValueWithNameSpace() {
	mockserver.Set("application", "key", "value")
	s.wait()
	val := GetStringValueWithNameSpace("application", "key", "defaultValue")
	s.Equal("value", val)
}

func (s *StartWithConfTestSuite) TestGetStringValue() {
	mockserver.Set("application", "key", "newvalue")
	s.wait()
	val := GetStringValue("key", "defaultValue")
	s.Equal("newvalue", val)
}

func (s *StartWithConfTestSuite) TestGetIntValue() {
	mockserver.Set("application", "intkey", "1")
	s.wait()
	val := GetIntValue("intkey", 0)
	s.Equal(1, val)
}
func (s *StartWithConfTestSuite) TestGetNameSpaceContent() {

	mockserver.Set("client.json", "content", `{"name":"apollo"}`)
	s.wait()

	val := GetNameSpaceContent("client.json", "{}")
	s.Equal(`{"name":"apollo"}`, val)
}

func (s *StartWithConfTestSuite) wait() {
	select {
	case <-s.changeEvent:
	case <-time.After(time.Second * 30):
	}
}
func startMockServer() {
	go func() {
		if err := mockserver.Run(); err != nil {
			log.Fatal(err)
		}
	}()
	// wait for mock server to run
	time.Sleep(time.Millisecond * 10)
}

func TestRunSuite(t *testing.T) {
	suite.Run(t, new(StartWithConfTestSuite))
}
func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	tearDown()
	os.Exit(code)
}
func setup() {
	startMockServer()
}
func tearDown() {
	mockserver.Close()
}
