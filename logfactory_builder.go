package glog

import (
	"fmt"
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
	CreateAppender(name string, writers map[string]io.Writer, cfg map[string]string) (IAppender, error)

	SetFilterFactory(name string, filter IFilterFactory)
	CreateFilter(name string, cfg map[string]string) (IFilter, error)

	SetWriterFactory(name string, writer IWriterFactory)
	CreateWriter(name string, cfg map[string]string) (io.Writer, error)

	Build(cfg *ConfigLogRoot) (ILogFactory, error)
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
		writerFactory:          make(map[string]IWriterFactory),
		elementFormaterFactory: make(map[string]IElementFormatterFactory),
		appenderFactory:        make(map[string]IAppenderFactory),
		filterFactory:          make(map[string]IFilterFactory),
	}
	for k, v := range ElementFormatterFactories {
		builder.elementFormaterFactory[k] = v
	}
	for k, v := range AppenderFactories {
		builder.appenderFactory[k] = v
	}
	for k, v := range FilterFactories {
		builder.filterFactory[k] = v
	}
	for k, v := range WriterFactories {
		builder.writerFactory[k] = v
	}
	return builder
}

func (this *logFactoryBuilder) SetWriterFactory(name string, writer IWriterFactory) {
	this.writerFactory[name] = writer
}

func (this *logFactoryBuilder) CreateWriter(name string, cfg map[string]string) (io.Writer, error) {
	factory := this.writerFactory[name]
	if factory == nil {
		return nil, fmt.Errorf("create writer %v failed, miss required factory", name)
	}
	return factory.NewWriter(cfg)
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

func (this *logFactoryBuilder) CreateFilter(name string, cfg map[string]string) (IFilter, error) {
	factory := this.filterFactory[name]
	if factory == nil {
		return nil, fmt.Errorf("create filter %v failed, miss required factory", name)
	}
	return factory.NewFilter(this, cfg)
}

func (this *logFactoryBuilder) CreateAppender(name string, writers map[string]io.Writer, cfg map[string]string) (IAppender, error) {
	factory := this.appenderFactory[name]
	if factory == nil {
		return nil, fmt.Errorf("create appender %v failed, miss required factory", name)
	}
	return factory.NewAppender(writers, cfg)
}

func (this *logFactoryBuilder) SetAppenderFactory(name string, appender IAppenderFactory) {
	this.appenderFactory[name] = appender
}

func (this *logFactoryBuilder) buildWriters(cfg map[string]map[string]string) (map[string]io.Writer, error) {
	writers := make(map[string]io.Writer)
	for writerName, writerCfg := range cfg {
		writerType, ok := writerCfg["writer"]
		if !ok {
			return nil, fmt.Errorf("create writer %v failed, miss required field %v", writerName, "writer")
		}
		writer, err := this.CreateWriter(writerType, writerCfg)
		if err != nil {
			return nil, NewComError(fmt.Sprintf("create writer %v failed", writerName), err)
		}
		writers[writerName] = writer
	}
	return writers, nil
}

type AppenderCtx struct {
	appender    IAppender
	prepareable []ILogPrepare
}

func (this *logFactoryBuilder) buildAppenders(cfg map[string]*ConfigAppender, writers map[string]io.Writer) (map[string]*AppenderCtx, error) {
	appenders := make(map[string]*AppenderCtx, len(cfg))
	for appenderName, appenderCfg := range cfg {
		appendCtx := &AppenderCtx{}
		appender, err := this.CreateAppender(appenderCfg.Appender, writers, appenderCfg.Params)
		if err != nil {
			return nil, NewComError("build appender error", err)
		}
		appendCtx.appender = appender
		for filterName, filterParam := range appenderCfg.Filters {
			filter, err := this.CreateFilter(filterName, filterParam)
			if err != nil {
				return nil, NewComError(fmt.Sprintf("build appender %v failed", appenderName), err)
			}
			appendCtx.appender.AddFilter(filter)
		}
		elements, format, err := this.layoutParser.LayoutParser([]byte(appenderCfg.Layout))
		if err != nil {
			return nil, NewComError(fmt.Sprintf("build appender %v failed, layout parse error", appenderName), err)
		}
		elementFormaters := make([]IElementFormatter, 0, len(elements))
		for i := 0; i < len(elements); i++ {
			ef, err := this.elementFormaterFactory[string(elements[i].Element)].NewElementFormatter(string(elements[i].Param))
			if err != nil {
				return nil, NewComError(fmt.Sprintf("build appender %v failed", appenderName), err)
			}
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
	return appenders, nil
}

func (this *logFactoryBuilder) Build(cfg *ConfigLogRoot) (ILogFactory, error) {
	writers, err := this.buildWriters(cfg.Writers)
	if err != nil {
		return nil, NewComError("build writer error", err)
	}
	appenders, err := this.buildAppenders(cfg.Appenders, writers)
	if err != nil {
		return nil, NewComError("build appenders error", err)
	}

	loggers := make(map[string]ILogger, len(cfg.Loggers))
	for name, lcfg := range cfg.Loggers {
		var firstAppender IAppender
		var prevAppender IAppender
		var preparable []ILogPrepare
		for _, appenderName := range lcfg.Appenders {
			nextAppender, ok := appenders[appenderName]
			if !ok {
				return nil, fmt.Errorf("miss required appender %v in logger %v", appenderName, name)
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

		nlog := NewLogger([]byte(name), nil, firstAppender, preparable)
		for filterName, filterParam := range lcfg.Filters {
			filter, err := this.CreateFilter(filterName, filterParam)
			if err != nil {
				return nil, NewComError(fmt.Sprintf("build logger %v error", name), err)
			}
			nlog.AddFilter(filter)
		}
		loggers[name] = nlog
	}
	return NewLogFactory(loggers), nil
}
