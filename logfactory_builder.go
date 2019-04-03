package glog

type LogConfig struct {
	appenders map[string]struct {
		params  map[string]string
		filters map[string]map[string]string
		layout  string
	}
	loggers map[string]struct {
		params    map[string]string
		filters   map[string]map[string]string
		appenders []string
	}
}

type ILogFactoryBuilder interface {
	GetAllLayoutParserNames() []string
	SetLayoutParser(name string, layoutParser ILayoutParser)
	GetLayoutParser(name string) ILayoutParser

	GetAllElementFormatterNames() []string
	SetElementFormatter(name string, elementFormatter IElementFormatter)
	GetElementFormater(name string) IElementFormatter

	GetAllAppenderNames() []string
	SetAppender(name string, appender IAppender)
	GetAppender(name string) IAppender

	GetAllLayoutFormaterNames() []string
	GetLayoutFormater(name string) ILayoutFormatter
	SetLayoutFormater(name string, formatter ILayoutFormatter)
}

type logFactoryBuilder struct {
	layoutParsers map[string]ILayoutParser
	layoutFormats map[string]ILayoutFormatter
	elements      map[string]IElementFormatter
	appenders     map[string]IAppender
}

func (this *logFactoryBuilder) GetAllLayoutFormaterNames() []string {
	panic("implement me")
}

func (this *logFactoryBuilder) GetLayoutFormater(name string) ILayoutFormatter {
	panic("implement me")
}

func (this *logFactoryBuilder) SetLayoutFormater(name string, formatter ILayoutFormatter) {
	panic("implement me")
}

func NewLogFactoryBuilder() ILogFactoryBuilder {
	builder := &logFactoryBuilder{}
	builder.layoutParsers[""] = FuncLayoutParser(DefaultLayoutParser)

	return builder
}

func (this *logFactoryBuilder) GetAllLayoutParserNames() []string {
	names := make([]string, 0, len(this.layoutParsers))
	for name, _ := range this.layoutParsers {
		names = append(names, name)
	}
	return names
}

func (this *logFactoryBuilder) SetLayoutParser(name string, layoutParser ILayoutParser) {
	this.layoutParsers[name] = layoutParser
}

func (this *logFactoryBuilder) GetLayoutParser(name string) ILayoutParser {
	return this.layoutParsers[name]
}

func (this *logFactoryBuilder) GetAllElementFormatterNames() []string {
	names := make([]string, 0, len(this.elements))
	for name, _ := range this.elements {
		names = append(names, name)
	}
	return names
}

func (this *logFactoryBuilder) SetElementFormatter(name string, elementFormatter IElementFormatter) {
	this.elements[name] = elementFormatter
}

func (this *logFactoryBuilder) GetElementFormater(name string) IElementFormatter {
	return this.elements[name]
}

func (this *logFactoryBuilder) GetAllAppenderNames() []string {
	names := make([]string, 0, len(this.appenders))
	for name, _ := range this.appenders {
		names = append(names, name)
	}
	return names
}

func (this *logFactoryBuilder) SetAppender(name string, appender IAppender) {
	this.appenders[name] = appender
}

func (this *logFactoryBuilder) GetAppender(name string) IAppender {
	return this.appenders[name]
}

func (this *logFactoryBuilder) Build(logcfg *LogConfig) ILogFactory {

}
