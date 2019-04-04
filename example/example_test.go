package example

import (
	"github.com/gosrv/glog"
	"testing"
	"time"
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

func Test4(t *testing.T) {
	factory, err := glog.LoadLogFactory("glog4.json")
	if err != nil {
		panic(err)
	}
	logger := factory.GetLogger("logger")
	// 输出go程id 日期时间 日志级别 fields 调用函数，文件行号
	logger.WithField("name", "eleven").Debug("hello glog")
	logger.WithFields(glog.LF{"name":"eleven", "age":18}).Info("hello glog")
}

func TestNormal(t *testing.T) {
	factory, err := glog.AutoLoadLogFactory()
	if err != nil {
		panic(err)
	}
	loggerDev := factory.GetLogger("loggerdev")
	loggerPub := factory.GetLogger("loggerpublish")
	// 开发中使用，输出所有级别到控制台和文件中,文件会每分钟分割一次，可配置
	for i:=0; i<120; i++ {
		loggerDev.Debug("hello dev")
		loggerDev.Info("hello dev")
		loggerDev.Error("hello dev")
		// 开发中使用，输出warn级别到文件中,文件会每分钟分割一次，可配置
		loggerPub.Debug("hello pub")
		loggerPub.Info("hello pub")
		loggerPub.Error("hello pub")
		time.Sleep(time.Second*2)
	}
}