package glog

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type SplitFileWriter struct {
	file    *os.File
	path    string
	ext     string
	spanidx int64
	span    int64
}

func NewSplitFileWriter(path string, ext string, span int64) (*SplitFileWriter, error) {
	ins := &SplitFileWriter{path: path, ext: ext, span: span}
	_, err := ins.prepareFile()
	if err != nil {
		return nil, NewComError("create sfile error", err)
	}
	return ins, nil
}

func (this *SplitFileWriter) Write(p []byte) (n int, err error) {
	writer, err := this.prepareFile()
	if err != nil {
		_, _ = os.Stderr.Write([]byte(err.Error() + "\n"))
	}
	return writer.Write(p)
}

func (this *SplitFileWriter) prepareFile() (io.Writer, error) {
	now := time.Now().Unix()
	curSpanIdx := now / this.span
	if curSpanIdx == this.spanidx && this.file != nil {
		return this.file, nil
	}
	this.spanidx = curSpanIdx
	if this.file != nil {
		_ = this.file.Close()
		this.file = nil
	}

	sufix := time.Unix(this.spanidx*this.span, 0).Format("2006-01-02 15:04:05")
	sufix = strings.Replace(sufix, ":", "", -1)
	fullName := this.path + sufix + "." + this.ext
	err := os.MkdirAll(filepath.Dir(fullName), 0644)
	if err != nil {
		return os.Stderr, err
	}

	this.file, err = os.OpenFile(fullName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return os.Stderr, err
	}

	return this.file, nil
}
