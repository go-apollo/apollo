# Gopkg.in/apollo

[![Build Status](https://cloud.drone.io/api/badges/go-apollo/apollo/status.svg)](https://cloud.drone.io/go-apollo/apollo)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-apollo/apollo)](https://goreportcard.com/report/github.com/go-apollo/apollo)

## Purpose

The goal of this project is to make the easiest way of using Ctrip apollo for golang applications. This project has been forked from [philchia/agollo](https://github.com/philchia/agollo) since 2018.8 but change a lot.

## Contributing

Fork -> Patch -> Push -> Pull Request

## Feature

- ✅ Multiple namespace support
- ✅ Fail tolerant
- ❎  YML to struct 
- ❎  Custom logger

## Required

**go 1.11** or later

## Installation

```sh
    go get -u gopkg.in/apollo.v0
    //if you use dep as your golang dep tool
    dep ensure -add  gopkg.in/apollo.v0
```

## Build
If you want build this project,should use go 1.11+
```
GO111MODULE=on; go mod download

```

## Usage

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
