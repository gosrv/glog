package glog

import (
	"fmt"
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
	NewLayoutFormatter(layout string, elementFormatters []IElementFormatter) ILayoutFormatter
}
type FuncLayoutFormatterFactory func(layout string, elementFormatters []IElementFormatter) ILayoutFormatter

func (this FuncLayoutFormatterFactory) NewLayoutFormatter(layout string, elementFormatters []IElementFormatter) ILayoutFormatter {
	return this(layout, elementFormatters)
}

type layoutFormatter struct {
	layout            string
	elementFormatters []IElementFormatter
}

func NewLayoutFormatter(layout string, elementFormatters []IElementFormatter) ILayoutFormatter {
	return &layoutFormatter{layout: layout, elementFormatters: elementFormatters}
}

func (this *layoutFormatter) LayoutFormat(param *LogParam) []byte {
	argLen := len(this.elementFormatters)
	args := make([]interface{}, argLen, argLen)
	for i := 0; i < argLen; i++ {
		args[i] = this.elementFormatters[i].ElementFormat(param)
	}
	return []byte(fmt.Sprintf(this.layout, args...))
}
