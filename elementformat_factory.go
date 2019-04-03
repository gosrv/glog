package glog

const (
	ElememtFormatBody  = "body"
	ElememtFormatDate  = "date"
	ElememtFormatLevel = "level"
)

func NewElememtFormatBody(param string) IElementFormatter {
	return FuncElementFormat(ElementFormatBody)
}

func NewElememtFormatDate(param string) IElementFormatter {
	return NewElementFormatDateTime(param)
}

func NewElememtFormatLevel(param string) IElementFormatter {
	return FuncElementFormat(ElementFormatLevel)
}

var ElementFormatterFactories = map[string]IElementFormatterFactory{
	ElememtFormatBody:  FuncElementFormatterFactory(NewElememtFormatBody),
	ElememtFormatDate:  FuncElementFormatterFactory(NewElememtFormatDate),
	ElememtFormatLevel: FuncElementFormatterFactory(NewElememtFormatLevel),
}
