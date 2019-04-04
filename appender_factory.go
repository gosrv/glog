package glog

import (
	"fmt"
	"io"
	"strconv"
)

const (
	AppenderChan   = "chan"
	AppenderWriter = "writer"

	AppenderWriterParamWriter = "writer"
)

func NewAppenderChan(writers map[string]io.Writer, params map[string]string) (IAppender, error) {
	pcap, ok := params["cap"]
	cap := 1024
	if ok {
		var err error
		cap, err = strconv.Atoi(pcap)
		if err != nil {
			return nil, fmt.Errorf("new appener chan error, param %v:%v parse error", "cap", pcap)
		}
	}

	return NewChanAppener(cap), nil
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
