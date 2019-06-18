//Copyright (c) 2017 Phil

package apollo

import (
	"log"
	"os"
	"testing"
	"time"

	// "time"

	"github.com/stretchr/testify/suite"

	"gopkg.in/apollo.v0/internal/mockserver"
)

type StartWithConfTestSuite struct {
	suite.Suite
	changeEvent <-chan *ChangeEvent
}

func (s *StartWithConfTestSuite) SetupSuite() {
	startMockServer()
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
	time.Sleep(5 * time.Second)
	mockserver.Close()
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

func (s *StartWithConfTestSuite) TestRunCacheSuite() {

	cs := new(CacheTestSuite)
	cs.SetT(s.T())
	s.Run("TestCache", cs.TestCache)
	s.Run("TestCacheDump", cs.TestCacheDump)
}

func (s *StartWithConfTestSuite) TestRunNotificationSuite() {

	ns := new(NotificationTestSuite)
	ns.SetT(s.T())
	s.Run("TestNotification", ns.TestNotification)
}
func (s *StartWithConfTestSuite) TestRunCommonSuite() {

	cs := new(CommonTestSuite)
	cs.SetT(s.T())
	s.Run("TestLocalIp", cs.TestLocalIp)
	s.Run("TestConfigURL", cs.TestConfigURL)
	s.Run("TestNotificationURL", cs.TestNotificationURL)
}
func (s *StartWithConfTestSuite) TestRunChangeSuite() {
	cs := new(ChangeTestSuite)
	cs.SetT(s.T())
	s.Run("TestCache", cs.TestChangeType)
	s.Run("TestMakeAddChange", cs.TestMakeAddChange)
	s.Run("TestMakeModifyChange", cs.TestMakeModifyChange)
	s.Run("TestMakeDeleteChange", cs.TestMakeDeleteChange)
}
func TestRunSuite(t *testing.T) {
	suite.Run(t, new(StartWithConfTestSuite))
}
