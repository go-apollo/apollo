//Copyright (c) 2017 Phil

package apollo

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Conf ...
type Conf struct {
	AppID      string   `yaml:"appId"`
	Cluster    string   `yaml:"cluster"`
	Namespaces []string `yaml:"namespaces,flow"`
	IP         string   `json:"ip"`
}

// NewConf create Conf from file
func NewConf(name string) (*Conf, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var ret Conf

	if err := yaml.NewDecoder(f).Decode(&ret); err != nil {
		return nil, err
	}

	return &ret, nil
}
