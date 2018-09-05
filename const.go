//Copyright (c) 2017 Phil

package apollo

import (
	"time"
)

const (
	defaultConfName  = "app.yml"
	defaultDumpFile  = ".apollo"
	defaultNamespace = "application"

	longPoolInterval      = time.Second * 2
	longPoolTimeout       = time.Second * 90
	queryTimeout          = time.Second * 2
	defaultNotificationID = -1
)
