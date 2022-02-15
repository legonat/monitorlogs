package erx

import (
	"errors"
	"fmt"
	"runtime"
)

func New(code error)(error){
	if code == nil {
		return code
	}
	_, fn, line, _ := runtime.Caller(1)
	text := fmt.Sprintf("Error File :\"%s) Line: %d Message (%s)", fn, line, code.Error())
	return errors.New(text)
}

func NewError(code int, message string)(error){

	_, fn, line, _ := runtime.Caller(1)
	text := fmt.Sprintf("Error File :\"%s) Line: %d Message (%s), Code (%d)", fn, line, message, code)
	return errors.New(text)
}