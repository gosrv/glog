package glog

import "fmt"

type ILayoutFormatter interface {
	LayoutFormat(param *LogParam) []byte
}
// {c}{level} {date:yy-MM-dd hh:mm:ss}	{body}	{fields:json} {file}:{fileline}{GOID}
type layoutFormatter struct {
	layout string
	elementFormatters []IElementFormatter
}

func NewLayoutFormatter(layout string, elementFormatters []IElementFormatter) *layoutFormatter {
	return &layoutFormatter{layout: layout, elementFormatters: elementFormatters}
}

func (this *layoutFormatter) LayoutFormat(param *LogParam) []byte {
	argLen := len(this.elementFormatters)
	args := make([]interface{}, argLen, argLen)
	for i:=0; i<argLen; i++ {
		args[i] = this.elementFormatters[i].ElementFormat(param)
	}
	return []byte(fmt.Sprintf(this.layout, args...))
}
