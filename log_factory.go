package glog

type ILogFactory interface {
	GetLogger(name string) ILogger
}

type logFactory struct {
	loggers map[string]ILogger
}

func NewLogFactory(loggers map[string]ILogger) *logFactory {
	return &logFactory{loggers: loggers}
}

func (this *logFactory) GetLogger(name string) ILogger {
	return this.loggers[name]
}
