package glog

import (
	"encoding/json"
	"io/ioutil"
)

const (
	autoLoadCfgFileName = "glog.json"
)

func AutoLoadLogFactory() (ILogFactory, error) {
	return LoadLogFactory(autoLoadCfgFileName)
}

func LoadLogFactory(path string) (ILogFactory, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfgroot := &ConfigLogRoot{}
	err = json.Unmarshal(data, cfgroot)
	if err != nil {
		return nil, err
	}
	builder := NewLogFactoryBuilder()
	return builder.Build(cfgroot), nil
}
