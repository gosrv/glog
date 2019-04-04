package example

import (
	"encoding/json"
	"fmt"
	"github.com/gosrv/glog"
	"testing"
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
      "layout": "{date:2006-01-02 15:04:05} [{level}] {body} {fields}"
    },
	"console" : {
      "params": {
      },
      "filters": {
        "level.limit": {"level": "debug"}
      },
      "layout": "{date:2006-01-02 15:04:05} [{level}] {body} {fields}"
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
    }
  }
}
`
func BenchmarkXX(b *testing.B) {
	cfgroot := &glog.ConfigLogRoot{}
	_ = json.Unmarshal([]byte(cfg), cfgroot)

	builder := glog.NewLogFactoryBuilder()
	factory := builder.Build(cfgroot)
	logger := factory.GetLogger("testlogger1")
	b.StartTimer()
	for i:=0; i<b.N; i++ {
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

	logger.WithFields(glog.LF{"abc":123,"rrr":666}).Debug("hello")
	fmt.Println("finish************")
}
