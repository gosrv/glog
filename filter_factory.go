package glog

import (
	"fmt"
	"strings"
)

const (
	FilterLevelLimit  = "level.limit"
	FilterLevelPass   = "level.pass"
	ParamFilterLevel  = "level"
	ParamFilterPass   = "pass"
	ParamFilterReject = "reject"
)

func NewLogLevelLimitFilter(builder ILogFactoryBuilder, params map[string]string) (IFilter, error) {
	plevel, ok := params[ParamFilterLevel]
	if !ok || len(plevel) == 0 {
		return nil, fmt.Errorf("new log level limit error, miss required param %v", ParamFilterLevel)
	}
	level, err := ParseLevel(plevel)
	if err != nil {
		return nil, NewComError("new log level limit error, parse level error", err)
	}
	return newLogLevelLimitFilter(level), nil
}

func NewLogLevelPassFilter(builder ILogFactoryBuilder, params map[string]string) (IFilter, error) {
	var passLevelStrs []string = nil
	if len(params[ParamFilterPass]) > 0 {
		passLevelStrs = strings.Split(params[ParamFilterPass], ",")
	}
	var rejectLevelStrs []string = nil
	if len(params[ParamFilterReject]) > 0 {
		rejectLevelStrs = strings.Split(params[ParamFilterReject], ",")
	}
	var passLevels []Level
	for _, pl := range passLevelStrs {
		l, e := ParseLevel(pl)
		if e != nil {
			return nil, NewComError("new log level pass filter error, parse level error", e)
		}
		passLevels = append(passLevels, l)
	}
	var rejectLevels []Level
	for _, pl := range rejectLevelStrs {
		l, e := ParseLevel(pl)
		if e != nil {
			return nil, NewComError("new log level pass filter error, parse level error", e)
		}
		rejectLevels = append(rejectLevels, l)
	}
	return newLogLevelPassFilter(passLevels, rejectLevels), nil
}

var FilterFactories = map[string]IFilterFactory{
	FilterLevelLimit: FuncFilterFactory(NewLogLevelLimitFilter),
	FilterLevelPass:  FuncFilterFactory(NewLogLevelPassFilter),
}
