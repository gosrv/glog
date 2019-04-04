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
	writerName := params[AppenderWriterParamWriter]
	if len(writerName) == 0 {
		return nil, fmt.Errorf("new appener writer error, miss required param %v", AppenderWriterParamWriter)
	}
	writer, ok := writers[writerName]
	if !ok {
		return nil, fmt.Errorf("new appener writer error, miss required writer %v", writerName)
	}

	return NewIOWriterAppender(writer), nil
}

var AppenderFactories = map[string]IAppenderFactory{
	AppenderChan:   FuncAppenderFactory(NewAppenderChan),
	AppenderWriter: FuncAppenderFactory(NewAppenderWriter),
}
