package glog

type ILogFactory interface {
	GetLogger(name string) IFieldLogger
}

type logFactory struct {
	loggers map[string]IFieldLogger
}

func NewLogFactory(loggers map[string]IFieldLogger) *logFactory {
	return &logFactory{loggers: loggers}
}

func (this *logFactory) GetLogger(name string) IFieldLogger {
	return this.loggers[name]
}
