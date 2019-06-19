//Copyright (c) 2017 Phil

package apollo

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	logger "gopkg.in/logger.v1"

	"gopkg.in/apollo.v0/internal/mockserver"
)

type StartWithConfTestSuite struct {
	suite.Suite
	changeEvent <-chan *ChangeEvent
}

func (s *StartWithConfTestSuite) SetupSuite() {
	startTestApollo(s)
	s.changeEvent = WatchUpdate()
	mockserver.Set("application", "getkeys", "value")
	s.wait()
}
func startTestApollo(s *StartWithConfTestSuite) {
	s.Error(Start())
	s.NoError(StartWithConfFile("./testdata/app.yml"))
	s.NoError(defaultClient.loadLocal(defaultDumpFile))
}
func (s *StartWithConfTestSuite) BeforeTest(suiteName, testName string) {
	// log.Println("suiteName: " + suiteName)
	// log.Println("testName: " + testName)
}
func (s *StartWithConfTestSuite) TearDownSuite() {
	s.NoError(Stop())
	os.Remove(defaultDumpFile)
}
func (s *StartWithConfTestSuite) TestLoadLocal() {
	err := defaultClient.loadLocal(defaultDumpFile)
	s.NoError(err)
}
func (s *StartWithConfTestSuite) TestLogger() {
	setDefaultLogger()
	s.IsType(&logger.Logger{}, log)
	setLogger()
	s.Equal(0, logger.GetOutputLevel())
	s.NotNil(log)
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
	val2 := defaultClient.GetStringValue("key", "defaultValue")
	s.Equal("newvalue", val)
	s.Equal("newvalue", val2)
}
func (s *StartWithConfTestSuite) TestListKeys() {
	s.NotEmpty(ListKeys(defaultNamespace))
}
func (s *StartWithConfTestSuite) TestGetIntValue() {
	mockserver.Set("application", "intkey", "1")
	s.wait()
	val := GetIntValue("intkey", 0)
	s.NotEqual(0, val)
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
		log.Fatal(mockserver.Run())
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
func setLogger() {
	mLog := logger.Std
	mLog.SetOutputLevel(0)
	SetLogger(mLog)
}
func setup() {
	setLogger()
	startMockServer()
}
func tearDown() {
	mockserver.Close()
}
