package glog

import (
	"fmt"
	"io"
	"os"
)

const (
	AppenderChan   = "chan"
	AppenderWriter = "writer"

	AppenderWriterParamWriter = "writer"
)

func NewAppenderChan(writers map[string]io.Writer, params map[string]string) (IAppender, error) {
	return NewIOWriterAppender(os.Stdout), nil
}

func NewAppenderWriter(writers map[string]io.Writer, params map[string]string) (IAppender, error) {
	writerType := params[AppenderWriterParamWriter]
	if len(writerType) == 0 {
		return nil, fmt.Errorf("new appener writer error, miss required param %v", AppenderWriterParamWriter)
	}

	return NewIOWriterAppender(writers[writerType]), nil
}

var AppenderFactories = map[string]IAppenderFactory{
	AppenderChan:   FuncAppenderFactory(NewAppenderChan),
	AppenderWriter: FuncAppenderFactory(NewAppenderWriter),
}
