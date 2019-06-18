//Copyright (c) 2017 Phil
package apollo

import (
	"github.com/stretchr/testify/suite"
)

type ChangeTestSuite struct {
	suite.Suite
}

func (s *ChangeTestSuite) TestChangeType() {
	var tps = []ChangeType{ADD, MODIFY, DELETE, ChangeType(-1)}
	var strs = []string{"ADD", "MODIFY", "DELETE", "UNKNOW"}
	for i, tp := range tps {
		s.True(tp.String() == strs[i])
	}
}

func (s *ChangeTestSuite) TestMakeDeleteChange() {
	change := makeDeleteChange("key", []byte("val"))
	s.True(change.ChangeType == DELETE)
	s.True(string(change.OldValue) == "val")
}

func (s *ChangeTestSuite) TestMakeModifyChange() {
	change := makeModifyChange("key", []byte("old"), []byte("new"))
	s.True(change.ChangeType == MODIFY)
	s.True(string(change.OldValue) == "old")
	s.True(string(change.NewValue) == "new")
}

func (s *ChangeTestSuite) TestMakeAddChange() {
	change := makeAddChange("key", []byte("value"))
	s.True(change.ChangeType == ADD)
	s.True(string(change.NewValue) == "value")

}

// func TestRunChangeSuite(t *testing.T) {
// 	suite.Run(t, new(CacheTestSuite))
// }
