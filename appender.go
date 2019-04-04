package glog

import (
	"io"
)

type IAppenderFactory interface {
	NewAppender(writers map[string]io.Writer, params map[string]string) (IAppender, error)
}
type FuncAppenderFactory func(writers map[string]io.Writer, params map[string]string) (IAppender, error)

func (this FuncAppenderFactory) NewAppender(writers map[string]io.Writer, params map[string]string) (IAppender, error) {
	return this(writers, params)
}

// 输出目标
type IAppender interface {
	Next(next IAppender)
	Write(param *LogParam)
	AddFilter(filter IFilter)
}

type ChanAppener struct {
	FilterBundle
	next IAppender
	log  chan LogParam
}

func NewChanAppener(cap int) IAppender {
	ins := &ChanAppener{
		log: make(chan LogParam, cap),
	}
	go func() {
		for {
			data, ok := <-ins.log
			if !ok {
				break
			}
			if ins.next != nil {
				ins.next.Write(&data)
			}
		}
	}()
	return ins
}

func (this *ChanAppener) Next(next IAppender) {
	this.next = next
}

func (this *ChanAppener) Write(param *LogParam) {
	this.log <- *param
}

type IOWriterAppender struct {
	FilterBundle
	fmter  ILayoutFormatter
	next   IAppender
	writer io.Writer
}

func (this *IOWriterAppender) SetLayoutFormat(formatter ILayoutFormatter) {
	this.fmter = formatter
}

func NewIOWriterAppender(writer io.Writer) IAppender {
	return &IOWriterAppender{
		writer: writer,
	}
}

func (this *IOWriterAppender) Next(next IAppender) {
	this.next = next
}

func (this *IOWriterAppender) Write(param *LogParam) {
	data := this.fmter.LayoutFormat(param)
	if this.writer != nil {
		_, _ = this.writer.Write(data)
	}
	if this.next != nil {
		this.next.Write(param)
	}
}
