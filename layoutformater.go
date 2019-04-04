package glog

import (
	"reflect"
	"sync"
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
	cache             []byte
	layout            [][]byte
	elementFormatters []IElementFormatter
	lock              sync.Mutex
}

func NewLayoutFormatter(layout [][]byte, elementFormatters []IElementFormatter) ILayoutFormatter {
	return &layoutFormatter{
		layout:            layout,
		elementFormatters: elementFormatters,
	}
}

func (this *layoutFormatter) LayoutFormat(param *LogParam) []byte {
	argLen := len(this.elementFormatters)
	formated := make([]byte, 0, 512)
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
