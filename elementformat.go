package glog

import "time"

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
	return time.Now().Format(this.layout)
}

func ElementFormatBody(param *LogParam) string {
	return param.body
}

func ElementFormatLevel(param *LogParam) string {
	return param.level.String()
}
