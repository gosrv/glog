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

func NewSplitFileWriter(path string, ext string, span int64) *SplitFileWriter {
	return &SplitFileWriter{path: path, ext: ext, span: span}
}

func (this *SplitFileWriter) Write(p []byte) (n int, err error) {
	return this.prepareFile().Write(p)
}

func (this *SplitFileWriter) prepareFile() io.Writer {
	now := time.Now().Unix()
	curSpanIdx := now % this.span
	if curSpanIdx == this.spanidx && this.file != nil {
		return this.file
	}
	this.spanidx = curSpanIdx
	if this.file != nil {
		_ = this.file.Close()
		this.file = nil
	}

	sufix := time.Unix(this.spanidx*this.span, 0).Format("2006-01-02 15:04:05")
	sufix = strings.Trim(sufix, ":")
	fullName := this.path + sufix + "." + this.ext
	err := os.MkdirAll(filepath.Dir(fullName), 0644)
	if err != nil {
		_, _ = os.Stderr.Write([]byte(err.Error()))
		return os.Stderr
	}

	this.file, err = os.OpenFile(fullName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		_, _ = os.Stderr.Write([]byte(err.Error()))
		return os.Stderr
	}

	return this.file
}
