package glog

import "strings"

const (
	ElememtBody   = "body"
	ElememtDate   = "date"
	ElememtLevel  = "level"
	ElementFields = "fields"
	ElementFile   = "file"
	ElementGOID   = "goid"
)

func NewElememtBody(param string) IElementFormatter {
	return FuncElementFormat(ElementFormatBody)
}

func NewElememtDate(param string) IElementFormatter {
	return NewElementFormatDateTime(param)
}

func NewElememtLevel(param string) IElementFormatter {
	return FuncElementFormat(ElementFormatLevel)
}

func NewElementFields(param string) IElementFormatter {
	return FuncElementFormat(ElementFormatFields)
}

func NewElementFile(param string) IElementFormatter {
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
	return NewElementFormatFile(sep, short)
}

func NewElementGOID(param string) IElementFormatter {
	return FuncElementFormat(ElementFormatGOID)
}

var ElementFormatterFactories = map[string]IElementFormatterFactory{
	ElememtBody:   FuncElementFormatterFactory(NewElememtBody),
	ElememtDate:   FuncElementFormatterFactory(NewElememtDate),
	ElememtLevel:  FuncElementFormatterFactory(NewElememtLevel),
	ElementFields: FuncElementFormatterFactory(NewElementFields),
	ElementFile:   FuncElementFormatterFactory(NewElementFile),
	ElementGOID:   FuncElementFormatterFactory(NewElementGOID),
}
