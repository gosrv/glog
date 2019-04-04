package glog

const (
	ElememtFormatBody  = "body"
	ElememtFormatDate  = "date"
	ElememtFormatLevel = "level"
	ElementFormatField = "fields"
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

func NewElementFormatFields(param string) IElementFormatter {
	return FuncElementFormat(ElementFormatFields)
}

var ElementFormatterFactories = map[string]IElementFormatterFactory{
	ElememtFormatBody:  FuncElementFormatterFactory(NewElememtFormatBody),
	ElememtFormatDate:  FuncElementFormatterFactory(NewElememtFormatDate),
	ElememtFormatLevel: FuncElementFormatterFactory(NewElememtFormatLevel),
	ElementFormatField: FuncElementFormatterFactory(NewElementFormatFields),
}
