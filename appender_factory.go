package glog

import (
	"io"
	"os"
)

const (
	AppenderChan   = "chan"
	AppenderWriter = "writer"

	AppenderWriterParamWriter = "writer"
)

func NewAppenderChan(writers map[string]io.Writer, params map[string]string) IAppender {
	return NewIOWriterAppender(os.Stdout)
}

func NewAppenderWriter(writers map[string]io.Writer, params map[string]string) IAppender {
	writerType := params[AppenderWriterParamWriter]
	if len(writerType) == 0 {
		panic("appender has no file name")
	}

	return NewIOWriterAppender(writers[writerType])
}

var AppenderFactories = map[string]IAppenderFactory{
	AppenderChan:   FuncAppenderFactory(NewAppenderChan),
	AppenderWriter: FuncAppenderFactory(NewAppenderWriter),
}
