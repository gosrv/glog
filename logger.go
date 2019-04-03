package glog

type ILogger interface {
	SetLevel(level Level)
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Print(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Fatal(format string, args ...interface{})
	Panic(format string, args ...interface{})
}

type IFieldLogger interface {
	ILogger
	IFilterBundle
	IFilter

	WithField(key string, value interface{}) IFieldLogger
	WithFields(fields map[string]interface{}) IFieldLogger
	CreateLoggerWithFields(fields map[string]interface{}) IFieldLogger
}

type LogParam struct {
	fixFields map[string]interface{}
	fields    map[string]interface{}
	level     Level
}

type logger struct {
	LogParam
	FilterBundle
	appender  IAppender
	formatter ILayoutFormatter
}

func NewLogger(fixFields map[string]interface{}, appender IAppender, formatter ILayoutFormatter) *logger {
	return &logger{
		appender:  appender,
		formatter: formatter,
		LogParam: LogParam{
			fixFields: fixFields,
		},
	}
}

func (this *logger) WithField(key string, value interface{}) IFieldLogger {
	if this.fields == nil {
		this.fields = map[string]interface{}{}
	}
	this.fields[key] = value
	return this
}

func (this *logger) WithFields(fields map[string]interface{}) IFieldLogger {
	if this.fields == nil {
		this.fields = map[string]interface{}{}
	}
	for k, v := range fields {
		this.fields[k] = v
	}
	return this
}

func (this *logger) CreateLoggerWithFields(fields map[string]interface{}) IFieldLogger {
	allFields := map[string]interface{}{}
	for k, v := range fields {
		allFields[k] = v
	}
	for k, v := range this.fixFields {
		allFields[k] = v
	}
	return NewLogger(allFields, this.appender, this.formatter)
}

func (this *logger) SetLevel(level Level) {
	this.level = level
}

func (this *logger) Debug(format string, args ...interface{}) {
	this.Log(DebugLevel, format, args...)
}

func (this *logger) Info(format string, args ...interface{}) {
	this.Log(InfoLevel, format, args...)
}

func (this *logger) Print(format string, args ...interface{}) {
	this.Log(TraceLevel, format, args...)
}

func (this *logger) Warn(format string, args ...interface{}) {
	this.Log(WarnLevel, format, args...)
}

func (this *logger) Error(format string, args ...interface{}) {
	this.Log(ErrorLevel, format, args...)
}

func (this *logger) Fatal(format string, args ...interface{}) {
	this.Log(FatalLevel, format, args...)
}

func (this *logger) Panic(format string, args ...interface{}) {
	this.Log(PanicLevel, format, args...)
}

func (this *logger) Log(level Level, format string, args ...interface{}) {
	if this.level < level {
		return
	}
	this.appender.Write(&this.LogParam)
	this.fields = nil
}
