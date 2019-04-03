package glog

import "os"

const (
	AppenderConsole = "console"
	AppenderFile    = "file"

	ParamLayoutFormatter = "formatter"
	ParamFileName        = "path"
)

func NewConsoleAppender(builder ILogFactoryBuilder, params map[string]string) IAppender {
	return NewIOWriterAppender(builder.GetLayoutFormater(params[ParamLayoutFormatter]), os.Stdout)
}

func NewFileAppender(builder ILogFactoryBuilder, params map[string]string) IAppender {
	fname := params[ParamFileName]
	if len(fname) == 0 {
		panic("appender has no file name")
	}
	file, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0)
	if err != nil {
		panic(err)
	}
	return NewIOWriterAppender(builder.GetLayoutFormater(params[ParamLayoutFormatter]), file)
}

var AppenderFactories = map[string]IAppenderFactory{
	AppenderConsole: FuncAppenderFactory(NewConsoleAppender),
	AppenderFile:    FuncAppenderFactory(NewFileAppender),
}
