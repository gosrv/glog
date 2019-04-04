package example

import (
	"encoding/json"
	"github.com/gosrv/glog"
	"testing"
)

var cfg = `
{
	"writers" : {
		"discard":{
			"writer":"discard"
		},
		"console":{
			"writer":"console"
		},
		"sfile":{
			"writer":"sfile",
			"path":"slog/log.log"
			"span":"1m"
		}
	},
  "appenders" : {
    "discard1" : {
	  "appender":"writer",
      "params": {
		"writer":"discard"
      },
      "filters": {
        "level.limit": {"level": "debug"}
      },
      "layout": "[goid:{goid}] {date:2006-01-02 15:04:05} [{level}] {body} {fields} {file::,short}"
    },
	"console1" : {
	  "appender":"writer",
      "params": {
		"writer":"console"
      },
      "filters": {
        "level.limit": {"level": "debug"}
      },
      "layout": "[goid:{goid}] {date:2006-01-02 15:04:05} [{level}] {body} {fields} {file}"
    },
	"file1" : {
	  "appender":"writer",
      "params": {
		"writer":"sfile"
      },
      "filters": {
        "level.limit": {"level": "debug"}
      },
      "layout": "[goid:{goid}] {date:2006-01-02 15:04:05} [{level}] {body} {fields} {file}"
    }
  },
  "loggers" : {
    "testlogger1" : {
      "params": {
      },
      "filters": {
        "level.pass": {"pass": "debug", "reject": "error"}
      },
      "appenders": ["console1", "file1"]
    },
	"testlogger2" : {
      "params": {
      },
      "filters": {
        "level.pass": {"pass": "debug", "reject": "error"}
      },
      "appenders": ["discard1"]
    }
  }
}
`

func BenchmarkXX(b *testing.B) {
	cfgroot := &glog.ConfigLogRoot{}
	_ = json.Unmarshal([]byte(cfg), cfgroot)

	builder := glog.NewLogFactoryBuilder()
	factory := builder.Build(cfgroot)
	logger := factory.GetLogger("testlogger2")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		logger.WithFields(glog.LF{"abc": 123, "rrr": 666}).Debug("hello")
	}
	b.StopTimer()
}

func TestXX(t *testing.T) {
	cfgroot := &glog.ConfigLogRoot{}
	_ = json.Unmarshal([]byte(cfg), cfgroot)

	builder := glog.NewLogFactoryBuilder()
	factory := builder.Build(cfgroot)
	logger := factory.GetLogger("testlogger1")
	clog := logger.CreateLoggerWithFields(glog.LF{"name": "eleven"})

	logger.WithFields(glog.LF{"abc": 123, "rrr": 666}).Debug("hello")
	clog.WithFields(glog.LF{"abc": 123, "rrr": 666}).Debug("wrold")
}
