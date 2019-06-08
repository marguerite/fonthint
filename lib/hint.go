package lib

import (
	"github.com/golang/freetype"
	"io/ioutil"
	"reflect"
)

// Hinting checks if a font has builtin hinting instructions by reading its "fpgm" table
// basically hinting instructions are stored in "cvt", "fpgm", "prep" table.
// fpgm are the key table because its the actual bytecode intepreter virtual machine
// https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6fpgm.html
// https://developer.apple.com/fonts/TrueType-Reference-Manual/RM03/Chap3.html
func Hinting(f string) (bool, error) {
	b, e := ioutil.ReadFile(f)
	if e != nil {
		return false, e
	}
	font, e := freetype.ParseFont(b)
	if e != nil {
		return false, e
	}
	p := reflect.ValueOf(font)
	v := reflect.Indirect(p)
	fpgm := v.FieldByName("fpgm")
	if fpgm.Len() > 0 {
		return true, nil
	}
	return false, nil
}