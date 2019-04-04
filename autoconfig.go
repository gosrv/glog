package glog

import (
	"encoding/json"
	"fmt"
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
		return nil, NewComError(fmt.Sprintf("read file %v error ", path), err)
	}
	cfgroot := &ConfigLogRoot{}
	err = json.Unmarshal(data, cfgroot)
	if err != nil {
		return nil, NewComError(fmt.Sprintf("unmarshal file %v error ", path), err)
	}
	builder := NewLogFactoryBuilder()
	return builder.Build(cfgroot)
}
