package glog

import "io"

type IAppenderFactory interface {
	NewAppender(builder ILogFactoryBuilder, params map[string]string) IAppender
}
type FuncAppenderFactory func(builder ILogFactoryBuilder, params map[string]string) IAppender

func (this FuncAppenderFactory) NewAppender(builder ILogFactoryBuilder, params map[string]string) IAppender {
	return this(builder, params)
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

func NewChanAppener(cap int) *ChanAppener {
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

func NewIOWriterAppender(fmter ILayoutFormatter, writer io.Writer) *IOWriterAppender {
	return &IOWriterAppender{
		fmter:  fmter,
		writer: writer,
	}
}

func (this *IOWriterAppender) Next(next IAppender) {
	this.next = next
}

func (this *IOWriterAppender) Write(param *LogParam) {
	data := this.fmter.LayoutFormat(param)
	_, _ = this.writer.Write(data)
	if this.next != nil {
		this.next.Write(param)
	}
}
