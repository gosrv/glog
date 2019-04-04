package glog

import (
	"io"
	"reflect"
)

type ConfigAppender struct {
	Writer   string
	Appender string
	Params   map[string]string
	Filters  map[string]map[string]string
	Layout   string
}

type ConfigLogger struct {
	Params    map[string]string
	Filters   map[string]map[string]string
	Appenders []string
}

type ConfigLogRoot struct {
	Writers   map[string]map[string]string
	Appenders map[string]*ConfigAppender
	Loggers   map[string]*ConfigLogger
}

type ILogFactoryBuilder interface {
	SetLayoutParser(layoutParser ILayoutParser)
	SetElementFormatterFactory(name string, formatter IElementFormatterFactory)
	SetLayoutFormatterFactory(layoutFormatterFactory ILayoutFormatterFactory)

	SetAppenderFactory(name string, appender IAppenderFactory)
	CreateAppender(name string, writers map[string]io.Writer, cfg map[string]string) IAppender

	SetFilterFactory(name string, filter IFilterFactory)
	CreateFilter(name string, cfg map[string]string) IFilter

	SetWriterFactory(name string, writer IWriterFactory)
	CreateWriter(name string, cfg map[string]string) io.Writer

	Build(cfg *ConfigLogRoot) ILogFactory
}

type logFactoryBuilder struct {
	layoutParser           ILayoutParser
	layoutFormatterFactory ILayoutFormatterFactory

	writerFactory          map[string]IWriterFactory
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

func (this *logFactoryBuilder) SetWriterFactory(name string, writer IWriterFactory) {
	this.writerFactory[name] = writer
}

func (this *logFactoryBuilder) CreateWriter(name string, cfg map[string]string) io.Writer {
	return this.writerFactory[name].NewWriter(cfg)
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

func (this *logFactoryBuilder) CreateAppender(name string, writers map[string]io.Writer, cfg map[string]string) IAppender {
	return this.appenderFactory[name].NewAppender(writers, cfg)
}

func (this *logFactoryBuilder) SetAppenderFactory(name string, appender IAppenderFactory) {
	this.appenderFactory[name] = appender
}

type AppenderCtx struct {
	appender    IAppender
	prepareable []ILogPrepare
}

func (this *logFactoryBuilder) Build(cfg *ConfigLogRoot) ILogFactory {
	writers := make(map[string]io.Writer)
	for writerName, writerCfg := range cfg.Writers {
		writer := this.CreateWriter(writerCfg["writer"], writerCfg)
		writers[writerName] = writer
	}

	appenders := make(map[string]*AppenderCtx, len(cfg.Appenders))
	for appenderName, appenderCfg := range cfg.Appenders {
		appendCtx := &AppenderCtx{}
		appendCtx.appender = this.CreateAppender(appenderCfg.Appender, writers, appenderCfg.Params)
		for filterName, filterParam := range appenderCfg.Filters {
			appendCtx.appender.AddFilter(this.CreateFilter(filterName, filterParam))
		}
		elements, format := this.layoutParser.LayoutParser([]byte(appenderCfg.Layout))
		elementFormaters := make([]IElementFormatter, 0, len(elements))
		for i := 0; i < len(elements); i++ {
			ef := this.elementFormaterFactory[string(elements[i].Element)].NewElementFormatter(string(elements[i].Param))
			elementFormaters = append(elementFormaters, ef)
			if reflect.TypeOf(ef).AssignableTo(ILogPrepareType) {
				appendCtx.prepareable = append(appendCtx.prepareable, ef.(ILogPrepare))
			}
		}
		layoutFormatter := this.layoutFormatterFactory.NewLayoutFormatter(format, elementFormaters)
		if reflect.TypeOf(appendCtx.appender).AssignableTo(ILayoutFormatterAwareType) {
			appendCtx.appender.(ILayoutFormatterAware).SetLayoutFormat(layoutFormatter)
		}
		appenders[appenderName] = appendCtx
	}

	loggers := make(map[string]ILogger, len(cfg.Loggers))
	for name, lcfg := range cfg.Loggers {
		var firstAppender IAppender
		var prevAppender IAppender
		var preparable []ILogPrepare
		for _, appenderName := range lcfg.Appenders {
			nextAppender, ok := appenders[appenderName]
			if !ok {
				panic("error")
			}
			preparable = append(preparable, nextAppender.prepareable...)
			if firstAppender == nil {
				firstAppender = nextAppender.appender
				prevAppender = nextAppender.appender
			} else {
				prevAppender.Next(nextAppender.appender)
				prevAppender = nextAppender.appender
			}
		}

		nlog := NewLogger(nil, firstAppender, preparable)
		for filterName, filterParam := range lcfg.Filters {
			nlog.AddFilter(this.CreateFilter(filterName, filterParam))
		}
		loggers[name] = nlog
	}
	return NewLogFactory(loggers)
}
