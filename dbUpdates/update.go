package dbUpdates

import (
	//Inbuild packages
	"fmt"
	"reflect"
	"unicode"
)

type Update struct{}

func (Update) Invoke(fileName string) error {
	var firstLetter rune
	for index, char := range fileName {
		if index == 0 {
			firstLetter = unicode.ToUpper(char)
			break
		}
	}
	fileName = fileName[1:]
	file := fmt.Sprintf("%s%s", string(firstLetter), fileName)
	values := reflect.ValueOf(Update{}).MethodByName(file).Call(nil)
	if len(values) > 0 && values[0].Type() == reflect.TypeOf((*error)(nil)).Elem() && values[0].Interface() != nil {
		return values[0].Interface().(error)
	}
	return nil
}
