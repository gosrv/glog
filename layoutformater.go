package glog

import (
	"reflect"
)

type ILayoutFormatter interface {
	LayoutFormat(param *LogParam) []byte
}

type ILayoutFormatterAware interface {
	SetLayoutFormat(formatter ILayoutFormatter)
}

var ILayoutFormatterAwareType = reflect.TypeOf((*ILayoutFormatterAware)(nil)).Elem()

type ILayoutFormatterFactory interface {
	NewLayoutFormatter(layout [][]byte, elementFormatters []IElementFormatter) ILayoutFormatter
}
type FuncLayoutFormatterFactory func(layout [][]byte, elementFormatters []IElementFormatter) ILayoutFormatter

func (this FuncLayoutFormatterFactory) NewLayoutFormatter(layout [][]byte, elementFormatters []IElementFormatter) ILayoutFormatter {
	return this(layout, elementFormatters)
}

type layoutFormatter struct {
	cache				[]byte
	layout            [][]byte
	elementFormatters []IElementFormatter
}

func NewLayoutFormatter(layout [][]byte, elementFormatters []IElementFormatter) ILayoutFormatter {
	llen := len(layout)
	if llen > 0 {
		last := layout[llen-1]
		if len(last) == 0 || last[len(last)-1] != '\n' {
			layout = append(layout, []byte{'\n'})
		}
	}

	return &layoutFormatter{
		cache:make([]byte, 0, 4096),
		layout: layout,
		elementFormatters: elementFormatters,
	}
}

func (this *layoutFormatter) LayoutFormat(param *LogParam) []byte {
	argLen := len(this.elementFormatters)
	formated := this.cache
	formated = append(formated, this.layout[0]...)
	for i := 0; i < argLen; i++ {
		arg := this.elementFormatters[i].ElementFormat(param)
		formated = append(formated, arg...)
		if i+1 < len(this.layout) {
			formated = append(formated, this.layout[i+1]...)
		}
	}
	return formated
}
