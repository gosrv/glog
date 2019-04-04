package glog

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	WriterStdOut   = "console"
	WriteFile      = "file"
	WriteSplitFile = "sfile"
	WriterDiscard  = "discard"

	ParamFilePath = "path"
	ParamFileSpan = "span"
)

type IWriterFactory interface {
	NewWriter(param map[string]string) io.Writer
}
type FuncWriterFactory func(param map[string]string) io.Writer

func (this FuncWriterFactory) NewWriter(param map[string]string) io.Writer {
	return this(param)
}

func NewWriterStd(param map[string]string) io.Writer {
	return os.Stdout
}

func NewWriterDiscard(param map[string]string) io.Writer {
	return ioutil.Discard
}

func NewWriterFile(param map[string]string) io.Writer {
	fname := param[ParamFilePath]
	if len(fname) == 0 {
		panic("appender has no file name")
	}
	_ = os.MkdirAll(filepath.Dir(fname), 0644)
	file, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	return file
}

func NewWriterSplitFile(param map[string]string) io.Writer {
	path := param[ParamFilePath]
	if len(path) == 0 {
		panic("appender has no file name")
	}
	lext := strings.LastIndex(path, ".")
	fname := path[:lext]
	fext := ""
	if lext < len(path) {
		fext = path[lext+1:]
	}
	span := param[ParamFileSpan]
	if len(span) < 2 {
		panic("error")
	}
	ispan, err := strconv.Atoi(span[:len(span)-1])
	if err != nil {
		panic(err)
	}
	switch span[len(span)-1] {
	case 's':
	case 'm':
		ispan *= 60
	case 'h':
		ispan *= 3600
	case 'd':
		ispan *= 24 * 3600
	default:
		panic("error")
	}

	return NewSplitFileWriter(fname, fext, int64(ispan))
}

var WriterFactories = map[string]IWriterFactory{
	WriterStdOut:   FuncWriterFactory(NewWriterStd),
	WriteFile:      FuncWriterFactory(NewWriterFile),
	WriteSplitFile: FuncWriterFactory(NewWriterSplitFile),
	WriterDiscard:  FuncWriterFactory(NewWriterDiscard),
}
