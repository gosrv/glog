package example

import (
	"github.com/gosrv/glog"
	"testing"
	"time"
)

func BenchmarkXX(b *testing.B) {
	factory, err := glog.AutoLoadLogFactory()
	if err != nil {
		panic(err)
	}

	logger := factory.GetLogger("testlogger2")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		logger.WithFields(glog.LF{"abc": 123, "rrr": 666}).Debug("hello")
	}
	b.StopTimer()
}

func TestXX(t *testing.T) {
	factory, err := glog.AutoLoadLogFactory()
	if err != nil {
		panic(err)
	}
	logger := factory.GetLogger("testlogger1")
	clog := logger.CreateLoggerWithFields(glog.LF{"name": "eleven"})
	for {
		logger.WithFields(glog.LF{"abc": 123, "rrr": 666}).Debug("hello")
		clog.WithFields(glog.LF{"abc": 123, "rrr": 666}).Debug("wrold")
		time.Sleep(10 * time.Second)
	}

}
