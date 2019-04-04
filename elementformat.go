package glog

import (
	"bytes"
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
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

func ElementFormatLogger(param *LogParam) []byte {
	return param.loggerName
}

func ElementFormatLevel(param *LogParam) []byte {
	data, _ := param.level.MarshalText()
	return data
}

func ElementFormatFields(param *LogParam) []byte {
	buf := make([]byte, 0, 512)
	for _, field := range param.fixFields {
		buf = append(buf, []byte(field.key)...)
		buf = append(buf, '=')
		buf = append(buf, []byte(fmt.Sprintf("%v", field.val))...)
		buf = append(buf, ' ')
	}

	for _, field := range param.fields {
		buf = append(buf, []byte(field.key)...)
		buf = append(buf, '=')
		buf = append(buf, []byte(fmt.Sprintf("%v", field.val))...)
		buf = append(buf, ' ')
	}
	return buf
}

func ElementFormatGOID(param *LogParam) []byte {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return []byte(strconv.Itoa(int(n)))
}

type ElementFormatFile struct {
	sep   string
	short bool
}

func NewElementFormatFile(sep string, short bool) *ElementFormatFile {
	if len(sep) == 0 {
		sep = ":"
	}
	return &ElementFormatFile{sep: sep, short: short}
}

func (this *ElementFormatFile) LogPrepare(param *LogParam) {
	_, okFile := param.prepare["file"]
	if okFile {
		return
	}
	filename, line, funcname := "???", 0, "???"
	pc, filename, line, ok := runtime.Caller(3)
	if ok {
		funcname = runtime.FuncForPC(pc).Name()
	}
	if this.short {
		funcname = filepath.Ext(funcname)            // .foo
		funcname = strings.TrimPrefix(funcname, ".") // foo
		filename = filepath.Base(filename)
	}
	param.prepare["file"] = filename
	param.prepare["line"] = strconv.Itoa(line)
	param.prepare["func"] = funcname
}

func (this *ElementFormatFile) ElementFormat(param *LogParam) []byte {
	name := param.prepare["file"]
	line := param.prepare["line"]
	funcn := param.prepare["func"]
	cache := make([]byte, 0, len(name)+len(line)+len(funcn)+len(this.sep))
	cache = append(cache, []byte(funcn)...)
	cache = append(cache, []byte(this.sep)...)
	cache = append(cache, []byte(line)...)
	cache = append(cache, []byte(this.sep)...)
	cache = append(cache, []byte(name)...)
	return cache
}
