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
	WithFields(fields LF) IFieldLogger
}

type LogParam struct {
	Fields   []LogField
	LogLevel Level
	Body     []byte
	Prepare  map[string]string
	LogName  []byte
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
	FilterBundle
	fields     []LogField
	logName    []byte
	appender   IAppender
	logPrepare []ILogPrepare
}

func NewLogger(name []byte, fields []LogField,
	appender IAppender, logPrepare []ILogPrepare) *logger {
	return &logger{
		appender:   appender,
		logName:    name,
		logPrepare: logPrepare,
		fields:     fields,
	}
}

func (this *logger) WithField(key string, value interface{}) IFieldLogger {
	return NewLogger(this.logName, append(this.fields, NewLogField(key, value)),
		this.appender, this.logPrepare)
}

func (this *logger) WithFields(fields LF) IFieldLogger {
	nfields := this.fields
	for fn, fe := range fields {
		nfields = append(nfields, NewLogField(fn, fe))
	}
	return NewLogger(this.logName, nfields,
		this.appender, this.logPrepare)
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
	logParam := LogParam{
		Fields:   this.fields,
		LogLevel: level,
		LogName:  this.logName,
	}
	if len(args) == 0 {
		logParam.Body = []byte(format)
	} else {
		logParam.Body = []byte(fmt.Sprintf(format, args...))
	}
	if !this.FilterBundle.IsLogPass(&logParam) {
		return
	}
	if len(this.logPrepare) > 0 {
		logParam.Prepare = make(map[string]string)
		for _, prepare := range this.logPrepare {
			prepare.LogPrepare(&logParam)
		}
	}

	this.appender.Write(&logParam)
}
