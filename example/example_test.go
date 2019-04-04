package example

import (
	"github.com/gosrv/glog"
	"testing"
)


func Test1(t *testing.T) {
	factory, err := glog.LoadLogFactory("glog1.json")
	if err != nil {
		panic(err)
	}
	logger := factory.GetLogger("logger")
	// 输出到控制台
	logger.Debug("hello glog")
}

func Test2(t *testing.T) {
	factory, err := glog.LoadLogFactory("glog2.json")
	if err != nil {
		panic(err)
	}
	logger := factory.GetLogger("logger")
	// 输出到控制台和文件
	logger.Debug("hello glog")
}

func Test3(t *testing.T) {
	factory, err := glog.LoadLogFactory("glog3.json")
	if err != nil {
		panic(err)
	}
	logger1 := factory.GetLogger("logger1")
	logger2 := factory.GetLogger("logger2")
	// 只输出info及以上日志
	logger1.Debug("log1 debug log")
	logger1.Info("log1 info log")
	// 只输出debug和info日志
	logger2.Error("log2 error log")
	logger2.Info("log2 info log")
}