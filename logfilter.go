package glog

type ILogFilter interface {
	IsLogPass(param *LogParam) bool
}

type logLevelLimitFilter struct {
	minLevel Level
}

func NewLogLevelLimitFilter(minLevel Level) *logLevelLimitFilter {
	return &logLevelLimitFilter{minLevel: minLevel}
}

func (this *logLevelLimitFilter) IsLogPass(param *LogParam) bool {
	return this.minLevel >= param.level
}

type logLevelPassFilter struct {
	passLevel []Level
	rejectLevel []Level
}

func NewLogLevelPassFilter(passLevel []Level, rejectLevel []Level) *logLevelPassFilter {
	return &logLevelPassFilter{passLevel: passLevel, rejectLevel: rejectLevel}
}

func (this *logLevelPassFilter) IsLogPass(param *LogParam) bool {
	for _,l := range this.rejectLevel {
		if l == param.level {
			return false
		}
	}
	if len(this.passLevel) == 0 {
		return true
	}
	for _, l := range this.passLevel {
		if l == param.level {
			return true
		}
	}
	return false
}

type IFilterBundle interface {
	AddFilter(filter ILogFilter)
	RmvFilter(filter ILogFilter)
}

type FilterBundle struct {
	filters []ILogFilter
}

func (this *FilterBundle)AddFilter(filter ILogFilter)  {
	this.filters = append(this.filters, filter)
}

func (this *FilterBundle) RmvFilter(filter ILogFilter) {
	for i,f := range this.filters {
		if f == filter {
			this.filters = append(this.filters[:i], this.filters[i+1:]...)
			return
		}
	}
}

func (this *FilterBundle) IsLogPass(param *LogParam) bool {
	for _, filter := range this.filters {
		if !filter.IsLogPass(param) {
			return false
		}
	}
	return true
}
