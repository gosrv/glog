package glog

type ILogFactory interface {
	GetLogger(name string) logger
}

type logFactory struct {
	loggers map[string]logger
}

func NewLogFactory(loggers map[string]logger, elements map[string]IElementFormatter, appenders map[string]IAppender) *logFactory {
	return &logFactory{loggers: loggers}
}

func (this *logFactory) GetLogger(name string) logger {
	return this.loggers[name]
}
