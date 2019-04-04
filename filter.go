package glog

type IFilter interface {
	IsLogPass(param *LogParam) bool
}
type IFilterFactory interface {
	NewFilter(builder ILogFactoryBuilder, params map[string]string) IFilter
}
type FuncFilterFactory func(builder ILogFactoryBuilder, params map[string]string) IFilter

func (this FuncFilterFactory) NewFilter(builder ILogFactoryBuilder, params map[string]string) IFilter {
	return this(builder, params)
}

type logLevelLimitFilter struct {
	minLevel Level
}

func newLogLevelLimitFilter(minLevel Level) *logLevelLimitFilter {
	return &logLevelLimitFilter{minLevel: minLevel}
}

func (this *logLevelLimitFilter) IsLogPass(param *LogParam) bool {
	return this.minLevel >= param.LogLevel
}

type logLevelPassFilter struct {
	passLevel   []Level
	rejectLevel []Level
}

func newLogLevelPassFilter(passLevel []Level, rejectLevel []Level) *logLevelPassFilter {
	return &logLevelPassFilter{passLevel: passLevel, rejectLevel: rejectLevel}
}

func (this *logLevelPassFilter) IsLogPass(param *LogParam) bool {
	for _, l := range this.rejectLevel {
		if l == param.LogLevel {
			return false
		}
	}
	if len(this.passLevel) == 0 {
		return true
	}
	for _, l := range this.passLevel {
		if l == param.LogLevel {
			return true
		}
	}
	return false
}

type IFilterBundle interface {
	AddFilter(filter IFilter)
	RmvFilter(filter IFilter)
}

type FilterBundle struct {
	filters []IFilter
}

func (this *FilterBundle) AddFilter(filter IFilter) {
	this.filters = append(this.filters, filter)
}

func (this *FilterBundle) RmvFilter(filter IFilter) {
	for i, f := range this.filters {
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
