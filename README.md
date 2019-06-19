# Gopkg.in/apollo

[![Build Status](https://cloud.drone.io/api/badges/go-apollo/apollo/status.svg)](https://cloud.drone.io/go-apollo/apollo)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-apollo/apollo)](https://goreportcard.com/report/github.com/go-apollo/apollo)
[![](https://godoc.org/gopkg.in/apollo.v0?status.svg)](http://godoc.org/gopkg.in/apollo.v0)
[![Coverage Status](https://coveralls.io/repos/github/go-apollo/apollo/badge.svg?branch=master)](https://coveralls.io/github/go-apollo/apollo?branch=master)
## Purpose

The goal of this project is to make the easiest way of using Ctrip apollo for golang applications. This project has been forked from [philchia/agollo](https://github.com/philchia/agollo) since 2018.8 but change a lot.

## Contributing

Fork -> Patch -> Push -> Pull Request

## Feature

- ✅ Multiple namespace support
- ✅ Fail tolerant
- ✅  Custom logger
- ❎  YML to struct 


## Required

**go 1.10** or later

## Build
If you want build this project,should use go 1.11+
```
GO111MODULE=on; go mod download

```

## Usage
### Installation
```bash
# go mod (only go 1.11+) or project in gopath(go 1.10 +)
go get -u gopkg.in/apollo.v0
# if you use dep as your golang dep tool (go 1.10)
dep ensure -add  gopkg.in/apollo.v0
```
### Set custom logger(Optional)
go-apoll use gopkg.in/logger.v1 as default logger provider.
Any logger implemented apollo.Logger interface can be use as apollo logger provider(such as [logrus](https://github.com/sirupsen/logrus)).
```golang
//Logger interface
type Logger interface {
	Warnf(format string, v ...interface{})
	Warn(v ...interface{})
	Errorf(format string, v ...interface{})
	Error(v ...interface{})
	Infof(format string, v ...interface{})
	Info(v ...interface{})
	Debugf(format string, v ...interface{})
	Debug(v ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}
```
set logrus as log provider
```golang
var log = logrus.New()
log.Formatter = new(logrus.JSONFormatter)
log.Formatter = new(logrus.TextFormatter)                     //default
log.Formatter.(*logrus.TextFormatter).DisableColors = true    // remove colors
log.Formatter.(*logrus.TextFormatter).DisableTimestamp = true // remove timestamp from test output
log.Level = logrus.TraceLevel
log.Out = os.Stdout

apollo.SetLogger(log)
```


### Start use default app.yml config file

```golang
    apollo.Start()
```

### Start use given config file path

```golang
    apollo.StartWithConfFile(name)
```

### Subscribe to updates

```golang
    events := apollo.WatchUpdate()
    changeEvent := <-events
    bytes, _ := json.Marshal(changeEvent)
    fmt.Println("event:", string(bytes))
```

### Get apollo values

```golang
    apollo.GetStringValue(Key, defaultValue)
    apollo.GetStringValueWithNameSapce(namespace, key, defaultValue)
    apollo.GetIntValue(Key, defaultValue)
    apollo.GetIntValueWithNameSapce(namespace, key, defaultValue)
```

### Get namespace file contents

```golang
    apollo.GetNameSpaceContent(namespace, defaultValue)
```

## License

apollo is released under MIT lecense
