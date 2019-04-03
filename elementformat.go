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

func (this FuncElementFormatterFactory) NewLayoutFormatter(param string) IElementFormatter {
	return this(param)
}

type DateTimeElementFormat struct {
}

func (this DateTimeElementFormat) ElementFormat(param *LogParam) string {
	return time.Now().String()
}

type ColorElementFormat struct {
}

func (this *ColorElementFormat) ElementFormat(param *LogParam) string {
	return ""
}

type LevelElementFormat struct {
}
