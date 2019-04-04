package example

import (
	"encoding/json"
	"github.com/gosrv/glog"
	"testing"
	"time"
)

var cfg = `
{
  "appenders" : {
    "discard" : {
      "params": {
      },
      "filters": {
        "level.limit": {"level": "debug"}
      },
      "layout": "[goid:{goid}] {date:2006-01-02 15:04:05} [{level}] {body} {fields} {file::,short}"
    },
	"console" : {
      "params": {
      },
      "filters": {
        "level.limit": {"level": "debug"}
      },
      "layout": "[goid:{goid}] date:2006-01-02 15:04:05} [{level}] {body} {fields} {file}"
    }
  },
  "loggers" : {
    "testlogger1" : {
      "params": {
      },
      "filters": {
        "level.pass": {"pass": "debug", "reject": "error"}
      },
      "appenders": ["console"]
    },
	"testlogger2" : {
      "params": {
      },
      "filters": {
        "level.pass": {"pass": "debug", "reject": "error"}
      },
      "appenders": ["discard"]
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

	go logger.WithFields(glog.LF{"abc": 123, "rrr": 666}).Debug("hello")
	go logger.WithFields(glog.LF{"abc": 123, "rrr": 666}).Debug("wrold")
	time.Sleep(time.Second)
}
