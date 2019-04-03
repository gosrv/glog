package glog

import "strings"

const (
	FilterLevelLimit  = "level.limit"
	FilterLevelPass   = "level.pass"
	ParamFilterLevel  = "level"
	ParamFilterPass   = "pass"
	ParamFilterReject = "reject"
)

func NewLogLevelLimitFilter(builder ILogFactoryBuilder, params map[string]string) IFilter {
	level, err := ParseLevel(params[ParamFilterLevel])
	if err != nil {
		panic(err)
	}
	return newLogLevelLimitFilter(level)
}

func NewLogLevelPassFilter(builder ILogFactoryBuilder, params map[string]string) IFilter {
	passLevelStrs := strings.Split(params[ParamFilterPass], ",")
	rejectLevelStrs := strings.Split(params[ParamFilterReject], ",")
	var passLevels []Level
	for _, pl := range passLevelStrs {
		l, e := ParseLevel(pl)
		if e != nil {
			panic(e)
		}
		passLevels = append(passLevels, l)
	}
	var rejectLevels []Level
	for _, pl := range rejectLevelStrs {
		l, e := ParseLevel(pl)
		if e != nil {
			panic(e)
		}
		rejectLevels = append(rejectLevels, l)
	}
	return newLogLevelPassFilter(passLevels, rejectLevels)
}

var FilterFactories = map[string]IFilterFactory{
	FilterLevelLimit: FuncFilterFactory(NewLogLevelLimitFilter),
	FilterLevelPass:  FuncFilterFactory(NewLogLevelPassFilter),
}
