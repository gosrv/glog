package glog

import "strings"

const (
	ElememtBody   = "body"
	ElememtDate   = "date"
	ElememtLevel  = "level"
	ElementFields = "fields"
	ElementFile   = "file"
	ElementGOID   = "goid"
	ElementLogger = "logger"
)

func NewElememtBody(param string) (IElementFormatter, error) {
	return FuncElementFormat(ElementFormatBody), nil
}

func NewElememtDate(param string) (IElementFormatter, error) {
	return NewElementFormatDateTime(param), nil
}

func NewElememtLevel(param string) (IElementFormatter, error) {
	return FuncElementFormat(ElementFormatLevel), nil
}

func NewElementFields(param string) (IElementFormatter, error) {
	return FuncElementFormat(ElementFormatFields), nil
}

func NewElementFile(param string) (IElementFormatter, error) {
	sp := strings.Split(param, ",")
	sep := ":"
	short := true
	switch len(sp) {
	case 1:
		sep = sp[0]
	case 2:
		sep = sp[0]
		short = sp[1] == "short"
	}
	return NewElementFormatFile(sep, short), nil
}

func NewElementGOID(param string) (IElementFormatter, error) {
	return FuncElementFormat(ElementFormatGOID), nil
}

func NewElementLogger(param string) (IElementFormatter, error) {
	return FuncElementFormat(ElementFormatLogger), nil
}

var ElementFormatterFactories = map[string]IElementFormatterFactory{
	ElememtBody:   FuncElementFormatterFactory(NewElememtBody),
	ElememtDate:   FuncElementFormatterFactory(NewElememtDate),
	ElememtLevel:  FuncElementFormatterFactory(NewElememtLevel),
	ElementFields: FuncElementFormatterFactory(NewElementFields),
	ElementFile:   FuncElementFormatterFactory(NewElementFile),
	ElementGOID:   FuncElementFormatterFactory(NewElementGOID),
	ElementLogger: FuncElementFormatterFactory(NewElementLogger),
}
