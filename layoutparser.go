package glog

import (
	"strings"
)

type LayoutElement struct {
	Element string
	Param   string
}

func NewLayoutElement(element string, param string) *LayoutElement {
	return &LayoutElement{Element: element, Param: param}
}

type ILayoutParser interface {
	LayoutParser(layout string) ([]*LayoutElement, error)
}
type FuncLayoutParser func(layout string) ([]*LayoutElement, error)

func (this FuncLayoutParser) LayoutParser(layout string) ([]*LayoutElement, error) {
	return this(layout)
}

type ILayoutParserFactory interface {
	NewLayoutParser(builder ILogFactoryBuilder, params map[string]string) ILayoutParser
}
type FuncLayoutParserFactory func(builder ILogFactoryBuilder, params map[string]string) ILayoutParser

func (this FuncLayoutParserFactory) NewLayoutParser(builder ILogFactoryBuilder, params map[string]string) ILayoutParser {
	return this(builder, params)
}

func readNextExpectChar(layout string, pos int, limit int, expect byte) int {
	for i := pos; i < limit; i++ {
		if layout[i] == expect {
			return i
		}
	}
	return limit
}

func readNextExpectCharNotInQuot(layout string, pos int, expect byte) int {
	quot := false
	for i := pos; i < len(layout); i++ {
		if layout[i] == '"' {
			quot = !quot
		}
		if quot {
			continue
		}
		if layout[i] == expect {
			return i
		}
	}
	return len(layout)
}

func readSpace(layout string, pos int) int {
	for i := pos; i < len(layout); i++ {
		switch layout[i] {
		case ' ':
		case '	':
		default:
			return i
		}
	}
	return len(layout)
}

func DefaultLayoutParser(layout string) ([]*LayoutElement, error) {
	pos := 0
	llen := len(layout)
	format := make([]byte, 0, llen)

	layoutElements := make([]*LayoutElement, 0, 4)
	layoutElements = append(layoutElements, nil)

	for pos < llen {
		// 读取开始符
		startPos := readNextExpectChar(layout, pos, len(layout), '{')
		// 读取结束符
		endPos := readNextExpectCharNotInQuot(layout, startPos+1, '}')
		// 之前的加入format
		format = append(format, layout[pos:startPos]...)
		format = append(format, "%v"...)
		colonPos := readNextExpectChar(layout, startPos, endPos, ':')
		element := strings.TrimSpace(layout[startPos+1 : colonPos])
		var value string
		if colonPos < endPos {
			value = strings.TrimSpace(layout[colonPos+1 : endPos])
			if len(value) > 0 && value[0] == '"' && value[len(value)-1] == '"' {
				value = strings.TrimSpace(value[1 : len(value)-1])
			}
		}
		layoutElements = append(layoutElements, NewLayoutElement(element, value))
		pos = endPos + 1
	}

	layoutElements[0] = NewLayoutElement(string(format), "")
	return layoutElements, nil
}
