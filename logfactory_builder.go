package glog

import "reflect"

type ConfigAppender struct {
	Params  map[string]string
	Filters map[string]map[string]string
	Layout  string
}

type ConfigLogger struct {
	Params    map[string]string
	Filters   map[string]map[string]string
	Appenders []string
}

type ConfigLogRoot struct {
	Appenders map[string]*ConfigAppender
	Loggers   map[string]*ConfigLogger
}

type ILogFactoryBuilder interface {
	SetLayoutParser(layoutParser ILayoutParser)
	SetElementFormatterFactory(name string, formatter IElementFormatterFactory)
	SetLayoutFormatterFactory(layoutFormatterFactory ILayoutFormatterFactory)

	SetAppenderFactory(name string, appender IAppenderFactory)
	CreateAppender(name string, cfg map[string]string) IAppender

	SetFilterFactory(name string, filter IFilterFactory)
	CreateFilter(name string, cfg map[string]string) IFilter

	Build(cfg *ConfigLogRoot) ILogFactory
}

type logFactoryBuilder struct {
	layoutParser           ILayoutParser
	layoutFormatterFactory ILayoutFormatterFactory

	elementFormaterFactory map[string]IElementFormatterFactory
	appenderFactory        map[string]IAppenderFactory
	filterFactory          map[string]IFilterFactory
}

func NewLogFactoryBuilder() ILogFactoryBuilder {
	builder := &logFactoryBuilder{
		layoutParser:           FuncLayoutParser(DefaultLayoutParser),
		layoutFormatterFactory: FuncLayoutFormatterFactory(NewLayoutFormatter),
		elementFormaterFactory: ElementFormatterFactories,
		appenderFactory:        AppenderFactories,
		filterFactory:          FilterFactories,
	}
	return builder
}

func (this *logFactoryBuilder) SetLayoutParser(layoutParser ILayoutParser) {
	this.layoutParser = layoutParser
}

func (this *logFactoryBuilder) SetElementFormatterFactory(name string, formatter IElementFormatterFactory) {
	this.elementFormaterFactory[name] = formatter
}

func (this *logFactoryBuilder) SetLayoutFormatterFactory(layoutFormatterFactory ILayoutFormatterFactory) {
	this.layoutFormatterFactory = layoutFormatterFactory
}

func (this *logFactoryBuilder) SetFilterFactory(name string, filter IFilterFactory) {
	this.filterFactory[name] = filter
}

func (this *logFactoryBuilder) CreateFilter(name string, cfg map[string]string) IFilter {
	return this.filterFactory[name].NewFilter(this, cfg)
}

func (this *logFactoryBuilder) CreateAppender(name string, cfg map[string]string) IAppender {
	return this.appenderFactory[name].NewAppender(this, cfg)
}

func (this *logFactoryBuilder) SetAppenderFactory(name string, appender IAppenderFactory) {
	this.appenderFactory[name] = appender
}

func (this *logFactoryBuilder) Build(cfg *ConfigLogRoot) ILogFactory {
	appenders := make(map[string]IAppender, len(cfg.Appenders))
	for appenderName, appenderCfg := range cfg.Appenders {
		appender := this.CreateAppender(appenderName, appenderCfg.Params)
		for filterName, filterParam := range appenderCfg.Filters {
			appender.AddFilter(this.CreateFilter(filterName, filterParam))
		}
		elements, format := this.layoutParser.LayoutParser([]byte(appenderCfg.Layout))
		elementFormaters := make([]IElementFormatter, 0, len(elements))
		for i := 0; i < len(elements); i++ {
			ef := this.elementFormaterFactory[string(elements[i].Element)].NewElementFormatter(string(elements[i].Param))
			elementFormaters = append(elementFormaters, ef)
		}
		layoutFormatter := this.layoutFormatterFactory.NewLayoutFormatter(format, elementFormaters)
		if reflect.TypeOf(appender).AssignableTo(ILayoutFormatterAwareType) {
			appender.(ILayoutFormatterAware).SetLayoutFormat(layoutFormatter)
		}
		appenders[appenderName] = appender
	}

	loggers := make(map[string]ILogger, len(cfg.Loggers))
	for name, lcfg := range cfg.Loggers {
		var firstAppender IAppender
		var prevAppender IAppender
		for _, appenderName := range lcfg.Appenders {
			nextAppender, ok := appenders[appenderName]
			if !ok {
				panic("error")
			}
			if firstAppender == nil {
				firstAppender = nextAppender
				prevAppender = nextAppender
			} else {
				prevAppender.Next(nextAppender)
				prevAppender = nextAppender
			}
		}

		nlog := NewLogger(nil, firstAppender)
		for filterName, filterParam := range lcfg.Filters {
			nlog.AddFilter(this.CreateFilter(filterName, filterParam))
		}
		loggers[name] = nlog
	}
	return NewLogFactory(loggers)
}
