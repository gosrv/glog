package glog

import "os"

const (
	AppenderConsole = "console"
	AppenderFile    = "file"
	AppenderDiscard    = "discard"
	ParamFileName   = "path"
)

func NewAppenderConsole(builder ILogFactoryBuilder, params map[string]string) IAppender {
	return NewIOWriterAppender(os.Stdout)
}

func NewAppenderDiscard(builder ILogFactoryBuilder, params map[string]string) IAppender {
	return NewIOWriterAppender(nil)
}

func NewAppenderFile(builder ILogFactoryBuilder, params map[string]string) IAppender {
	fname := params[ParamFileName]
	if len(fname) == 0 {
		panic("appender has no file name")
	}
	file, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0)
	if err != nil {
		panic(err)
	}
	return NewIOWriterAppender(file)
}

var AppenderFactories = map[string]IAppenderFactory{
	AppenderConsole: FuncAppenderFactory(NewAppenderConsole),
	AppenderFile:    FuncAppenderFactory(NewAppenderFile),
	AppenderDiscard:    FuncAppenderFactory(NewAppenderDiscard),
}
