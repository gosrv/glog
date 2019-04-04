package glog

import (
	"fmt"
	"strings"
)

type Level uint32

func ToLevel(lvl string) (Level, error) {
	switch strings.ToLower(lvl) {
	case "panic":
		return PanicLevel, nil
	case "fatal":
		return FatalLevel, nil
	case "error":
		return ErrorLevel, nil
	case "warn":
		return WarnLevel, nil
	case "info":
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	case "trace":
		return TraceLevel, nil
	}

	return 0, fmt.Errorf("unknown level %v", lvl)
}

var btrace = []byte("trace")
var bdebug = []byte("debug")
var binfo = []byte("info")
var bwarn = []byte("warn")
var berror = []byte("error")
var bfatal = []byte("fatal")
var bpanic = []byte("panic")

func (level Level) MarshalText() ([]byte, error) {
	switch level {
	case TraceLevel:
		return btrace, nil
	case DebugLevel:
		return bdebug, nil
	case InfoLevel:
		return binfo, nil
	case WarnLevel:
		return bwarn, nil
	case ErrorLevel:
		return berror, nil
	case FatalLevel:
		return bfatal, nil
	case PanicLevel:
		return bpanic, nil
	}

	return nil, fmt.Errorf("unknown level %v", level)
}

var strace = "trace"
var sdebug = "debug"
var sinfo = "info"
var swarn = "warn"
var serror = "error"
var sfatal = "fatal"
var spanic = "panic"
var sunknown = "unknown"

func (level Level) String() string {
	switch level {
	case TraceLevel:
		return strace
	case DebugLevel:
		return sdebug
	case InfoLevel:
		return sinfo
	case WarnLevel:
		return swarn
	case ErrorLevel:
		return serror
	case FatalLevel:
		return sfatal
	case PanicLevel:
		return spanic
	}

	return sunknown
}

// A constant exposing all logging levels
var AllLevels = []Level{
	PanicLevel,
	FatalLevel,
	ErrorLevel,
	WarnLevel,
	InfoLevel,
	DebugLevel,
	TraceLevel,
}

const (
	PanicLevel Level = 1
	FatalLevel Level = 2
	ErrorLevel Level = 3
	WarnLevel  Level = 4
	InfoLevel  Level = 5
	DebugLevel Level = 6
	TraceLevel Level = 7
)

func ParseLevel(lvl string) (Level, error) {
	switch strings.ToLower(lvl) {
	case "panic":
		return PanicLevel, nil
	case "fatal":
		return FatalLevel, nil
	case "error":
		return ErrorLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "info":
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	case "trace":
		return TraceLevel, nil
	}

	var l Level
	return l, fmt.Errorf("not a valid Level: %q", lvl)
}
