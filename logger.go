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
	fixFields []LogField
	fields []LogField
	level  Level
	body   []byte
}

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
	appender  IAppender
}

func NewLogger(fixFields map[string]interface{}, appender IAppender) *logger {
	l := &logger{
		appender:  appender,
		LogParam:  LogParam{},
	}
	for fn, fe := range fixFields {
		l.fixFields = append(l.fixFields, NewLogField(fn, fe))
	}
	return l
}

func (this *logger) WithField(key string, value interface{}) IFieldLogger {
	this.fields = append(this.fields, NewLogField(key, value))
	return this
}

func (this *logger) WithFields(fields LF) IFieldLogger {
	for fn, fe := range fields {
		this.fields = append(this.fields, NewLogField(fn, fe))
	}
	return this
}

func (this *logger) CreateLoggerWithFields(fields LF) IFieldLogger {
	nfields := make(map[string]interface{}, len(fields) + len(this.fixFields))
	for _, fe := range this.fixFields {
		nfields[fe.key] = fe.val
	}
	for fn, fe := range fields {
		nfields[fn] = fe
	}
	return NewLogger(nfields, this.appender)
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
	if len(args) > 0 {
		this.LogParam.body = []byte(format)
	} else {
		this.LogParam.body = []byte(fmt.Sprintf(format, args...))
	}
	this.LogParam.level = level
	if !this.FilterBundle.IsLogPass(&this.LogParam) {
		return
	}
	this.appender.Write(&this.LogParam)
	this.fields = nil
}
