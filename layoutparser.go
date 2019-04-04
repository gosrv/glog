package glog

import (
	"bytes"
)

type LayoutElement struct {
	Element []byte
	Param   []byte
}

func NewLayoutElement(element []byte, param []byte) *LayoutElement {
	return &LayoutElement{Element: element, Param: param}
}

type ILayoutParser interface {
	LayoutParser(layout []byte) ([]*LayoutElement, [][]byte)
}
type FuncLayoutParser func(layout []byte) ([]*LayoutElement, [][]byte)

func (this FuncLayoutParser) LayoutParser(layout []byte) ([]*LayoutElement, [][]byte) {
	return this(layout)
}

func readNextExpectChar(layout []byte, pos int, limit int, expect byte) int {
	for i := pos; i < limit; i++ {
		if layout[i] == expect {
			return i
		}
	}
	return limit
}

func readNextExpectCharNotInQuot(layout []byte, pos int, expect byte) int {
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

func readSpace(layout []byte, pos int) int {
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

func DefaultLayoutParser(layout []byte) ([]*LayoutElement, [][]byte) {
	pos := 0
	llen := len(layout)
	format := make([][]byte, 0, llen)
	layoutElements := make([]*LayoutElement, 0, 16)

	for {
		// 读取开始符
		startPos := readNextExpectChar(layout, pos, len(layout), '{')
		// 读取结束符
		endPos := readNextExpectCharNotInQuot(layout, startPos+1, '}')
		// 之前的加入format
		if startPos == endPos {
			format = append(format, append(layout[pos:startPos], '\n'))
			break
		} else {
			format = append(format, layout[pos:startPos])
		}
		colonPos := readNextExpectChar(layout, startPos, endPos, ':')
		element := bytes.TrimSpace(layout[startPos+1 : colonPos])
		var value []byte
		if colonPos < endPos {
			value = bytes.TrimSpace(layout[colonPos+1 : endPos])
			if len(value) > 0 && value[0] == '"' && value[len(value)-1] == '"' {
				value = bytes.TrimSpace(value[1 : len(value)-1])
			}
		}
		layoutElements = append(layoutElements, NewLayoutElement(element, value))
		pos = endPos + 1
	}

	return layoutElements, format
}
