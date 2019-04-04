package glog

import (
	"errors"
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
	NewWriter(param map[string]string) (io.Writer, error)
}
type FuncWriterFactory func(param map[string]string) (io.Writer, error)

func (this FuncWriterFactory) NewWriter(param map[string]string) (io.Writer, error) {
	return this(param)
}

func NewWriterStd(param map[string]string) (io.Writer, error) {
	return os.Stdout, nil
}

func NewWriterDiscard(param map[string]string) (io.Writer, error) {
	return ioutil.Discard, nil
}

func NewWriterFile(param map[string]string) (io.Writer, error) {
	fname := param[ParamFilePath]
	if len(fname) == 0 {
		return nil, errors.New("create file writer failed, miss required param " + ParamFilePath)
	}
	err := os.MkdirAll(filepath.Dir(fname), 0644)
	if err != nil {
		return nil, NewComError("create file writer failed, mkdir error", err)
	}
	file, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, NewComError("create file writer failed, open file error", err)
	}

	return file, nil
}

func NewWriterSplitFile(param map[string]string) (io.Writer, error) {
	path := param[ParamFilePath]
	if len(path) == 0 {
		return nil, errors.New("create sfile writer failed, miss required param " + ParamFilePath)
	}
	lext := strings.LastIndex(path, ".")
	fname := path[:lext]
	fext := ""
	if lext < len(path) {
		fext = path[lext+1:]
	}
	span := param[ParamFileSpan]
	if len(span) < 2 {
		return nil, errors.New("create sfile writer failed, span format error " + span)
	}
	ispan, err := strconv.Atoi(span[:len(span)-1])
	if err != nil {
		return nil, errors.New("create sfile writer failed, span format error " + span)
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
		return nil, errors.New("create sfile writer failed, span format error " + span)
	}

	return NewSplitFileWriter(fname, fext, int64(ispan))
}

var WriterFactories = map[string]IWriterFactory{
	WriterStdOut:   FuncWriterFactory(NewWriterStd),
	WriteFile:      FuncWriterFactory(NewWriterFile),
	WriteSplitFile: FuncWriterFactory(NewWriterSplitFile),
	WriterDiscard:  FuncWriterFactory(NewWriterDiscard),
}
