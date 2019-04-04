package glog

import (
	"fmt"
	"time"
)

// 输出格式化
type IElementFormatter interface {
	ElementFormat(param *LogParam) []byte
}
type FuncElementFormat func(param *LogParam) []byte

func (this FuncElementFormat) ElementFormat(param *LogParam) []byte {
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

func (this *ElementFormatDateTime) ElementFormat(param *LogParam) []byte {
	if len(this.layout) > 0 {
		return []byte(time.Now().Format(this.layout))
	} else {
		data, _ := time.Now().MarshalBinary()
		return data
	}
}

func ElementFormatBody(param *LogParam) []byte {
	return param.body
}

func ElementFormatLevel(param *LogParam) []byte {
	data, _ := param.level.MarshalText()
	return data
}

func ElementFormatFields(param *LogParam) []byte {
	buf := make([]byte, 0, 512)
	for _, field := range param.fields {
		buf = append(buf, []byte(field.key)...)
		buf = append(buf, ':')
		buf = append(buf, []byte(fmt.Sprintf("%v", field.val))...)
		buf = append(buf, ' ')
	}
	return buf
}