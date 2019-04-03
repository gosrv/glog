package glog

import "os"

func NewConsoleAppender(builder ILogFactoryBuilder, params map[string]string) IAppender {
	return NewIOWriterAppender(builder.GetLayoutFormater(""), os.Stdout)
}