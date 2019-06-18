//Copyright (c) 2017 Phil

package apollo

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type NotificationTestSuite struct {
	suite.Suite
}

func (s *NotificationTestSuite) TestNotification() {
	repo := new(notificationRepo)

	repo.setNotificationID("namespace", 1)
	id, ok := repo.getNotificationID("namespace")
	s.True(ok)
	s.Equal(1, id)
	repo.setNotificationID("namespace", 2)
	id2, ok2 := repo.getNotificationID("namespace")
	s.True(ok2)
	s.Equal(2, id2)

	id3, ok3 := repo.getNotificationID("null")

	s.False(ok3)
	s.Equal(defaultNotificationID, id3)

	str := repo.toString()
	s.NotEmpty(str)
}
func TestRunNotificationSuite(t *testing.T) {
	suite.Run(t, new(NotificationTestSuite))
}
