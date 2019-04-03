package glog

import (
	"fmt"
	"testing"
)

func TestLayoutFormat(t *testing.T) {
	lt := "{ c}{level } { date  	: \"  yy-MM-dd hh:mm:ss \" }	{ body}	{fields :json} {file }:{fileline }{GOID}"
	val, _ := DefaultLayoutParser(lt)
	fmt.Println(val)
}