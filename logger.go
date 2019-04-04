package glog

import (
	"fmt"
	"reflect"
)

type ISharable interface {
	IsSharable() bool
}

var ISharableType = reflect.TypeOf((*ISharable)(nil)).Elem()

type ILogComponentInit interface {
	LogComponentInit(builder ILogFactoryBuilder, params map[string]string) error
}

type ILogger interface {
	WithField(key string, value interface{}) IFieldLogger
	WithFields(fields LF) IFieldLogger
	CreateLoggerWithFields(fields LF) IFieldLogger

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
}

type LogParam struct {
	FixFields []LogField
	Fields    []LogField
	LogLevel  Level
	Body      []byte
	Prepare   map[string]string
	LogName   []byte
}

type ILogPrepare interface {
	LogPrepare(param *LogParam)
}

var ILogPrepareType = reflect.TypeOf((*ILogPrepare)(nil)).Elem()

type LogField struct {
	key string
	val interface{}
}

type LF map[string]interface{}

func NewLogField(key string, val interface{}) LogField {
	return LogField{key: key, val: val}
}

type logger struct {
	LogParam
	FilterBundle
	appender   IAppender
	logPrepare []ILogPrepare
}

func NewLogger(name []byte, fixFields map[string]interface{},
	appender IAppender, logPrepare []ILogPrepare) *logger {
	l := &logger{
		appender: appender,
		LogParam: LogParam{
			Prepare: make(map[string]string),
			LogName: name,
		},
		logPrepare: logPrepare,
	}
	for fn, fe := range fixFields {
		l.FixFields = append(l.FixFields, NewLogField(fn, fe))
	}
	return l
}

func (this *logger) WithField(key string, value interface{}) IFieldLogger {
	this.Fields = append(this.Fields, NewLogField(key, value))
	return this
}

func (this *logger) WithFields(fields LF) IFieldLogger {
	for fn, fe := range fields {
		this.Fields = append(this.Fields, NewLogField(fn, fe))
	}
	return this
}

func (this *logger) CreateLoggerWithFields(fields LF) IFieldLogger {
	nfields := make(map[string]interface{}, len(fields)+len(this.FixFields))
	for _, fe := range this.FixFields {
		nfields[fe.key] = fe.val
	}
	for fn, fe := range fields {
		nfields[fn] = fe
	}
	return NewLogger(this.LogName, nfields, this.appender, this.logPrepare)
}

func (this *logger) Debug(format string, args ...interface{}) {
	this.log(DebugLevel, format, args...)
}

func (this *logger) Info(format string, args ...interface{}) {
	this.log(InfoLevel, format, args...)
}

func (this *logger) Print(format string, args ...interface{}) {
	this.log(TraceLevel, format, args...)
}

func (this *logger) Warn(format string, args ...interface{}) {
	this.log(WarnLevel, format, args...)
}

func (this *logger) Error(format string, args ...interface{}) {
	this.log(ErrorLevel, format, args...)
}

func (this *logger) Fatal(format string, args ...interface{}) {
	this.log(FatalLevel, format, args...)
}

func (this *logger) Panic(format string, args ...interface{}) {
	this.log(PanicLevel, format, args...)
}

func (this *logger) log(level Level, format string, args ...interface{}) {
	if len(args) == 0 {
		this.LogParam.Body = []byte(format)
	} else {
		this.LogParam.Body = []byte(fmt.Sprintf(format, args...))
	}
	this.LogParam.LogLevel = level
	if !this.FilterBundle.IsLogPass(&this.LogParam) {
		return
	}
	for _, prepare := range this.logPrepare {
		prepare.LogPrepare(&this.LogParam)
	}
	this.appender.Write(&this.LogParam)
	this.Fields = nil
	if len(this.LogParam.Prepare) > 0 {
		this.LogParam.Prepare = make(map[string]string)
	}
}
