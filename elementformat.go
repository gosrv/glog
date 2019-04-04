package glog

import (
	"fmt"
	"time"
)

// 输出格式化
type IElementFormatter interface {
	ElementFormat(param *LogParam) string
}
type FuncElementFormat func(param *LogParam) string

func (this FuncElementFormat) ElementFormat(param *LogParam) string {
	return this(param)
}

type IElementFormatterFactory interface {
	NewElementFormatter(params string) IElementFormatter
}
type FuncElementFormatterFactory func(param string) IElementFormatter

func (this FuncElementFormatterFactory) NewElementFormatter(param string) IElementFormatter {
	return this(param)
}

type ElementFormatDateTime struct {
	layout string
}

func NewElementFormatDateTime(layout string) IElementFormatter {
	return &ElementFormatDateTime{layout: layout}
}

func (this *ElementFormatDateTime) ElementFormat(param *LogParam) string {
	if len(this.layout) > 0 {
		return time.Now().Format(this.layout)
	} else {
		return time.Now().String()
	}
}

func ElementFormatBody(param *LogParam) string {
	return param.body
}

func ElementFormatLevel(param *LogParam) string {
	return param.level.String()
}

func ElementFormatFields(param *LogParam) string {
	buf := make([]byte, 0, 512)
	for _, field := range param.fields {
		buf = append(buf, []byte(field.key)...)
		buf = append(buf, ':')
		buf = append(buf, []byte(fmt.Sprintf("%v", field.val))...)
		buf = append(buf, ' ')
	}
	return string(buf)
}