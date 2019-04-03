package example

import (
	"encoding/json"
	"github.com/gosrv/glog"
	"testing"
)

var cfg = `
{
  "appenders" : {
    "console" : {
      "params": {
      },
      "filters": {
        "level.limit": {"level": "debug"}
      },
      "layout": "{date:2006-01-02 15:04:05} [{level}] {body}"
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

func TestXX(t *testing.T) {
	cfgroot := &glog.ConfigLogRoot{}
	_ = json.Unmarshal([]byte(cfg), cfgroot)

	builder := glog.NewLogFactoryBuilder()
	factory := builder.Build(cfgroot)
	logger := factory.GetLogger("testlogger1")
	logger.Debug("hello")
}
