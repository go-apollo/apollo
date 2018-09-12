//Copyright (c) 2017 Phil

//Package apollo ctrip apollo go client
package apollo

import log "gopkg.in/logger.v1"

var (
	defaultClient *Client
)

// Start apollo
func Start() error {
	return StartWithConfFile(defaultConfName)
}

// StartWithConfFile run apollo with conf file
func StartWithConfFile(name string) error {
	log.Debugf("StartWithConfFile run apollo with conf file name: %s", name)
	conf, err := NewConf(name)
	if err != nil {
		return err
	}
	return StartWithConf(conf)
}

// StartWithConf run apollo with Conf
func StartWithConf(conf *Conf) error {
	defaultClient = NewClient(conf)

	return defaultClient.Start()
}

// Stop sync config
func Stop() error {
	return defaultClient.Stop()
}

// WatchUpdate get all updates
func WatchUpdate() <-chan *ChangeEvent {
	return defaultClient.WatchUpdate()
}

// GetStringValueWithNameSpace get value from given namespace
func GetStringValueWithNameSpace(namespace, key, defaultValue string) string {
	return defaultClient.GetStringValueWithNameSpace(namespace, key, defaultValue)
}

// GetStringValue from default namespace
func GetStringValue(key, defaultValue string) string {
	return GetStringValueWithNameSpace(defaultNamespace, key, defaultValue)
}

// GetNameSpaceContent get contents of namespace
func GetNameSpaceContent(namespace, defaultValue string) string {
	return defaultClient.GetNameSpaceContent(namespace, defaultValue)
}
