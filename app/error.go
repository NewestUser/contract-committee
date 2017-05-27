package app

import "fmt"

type NotFoundError struct {
	msg string
}

func (e *NotFoundError) Error() string {
	return e.msg
}

func NotFoundErr(format string, a ...interface{}) error {
	return &NotFoundError{msg:fmt.Sprintf(format, a...)}
}

type BadRequestError struct {
	msg string
}

func (e *BadRequestError) Error() string {
	return e.msg
}

func BadReqErr(format string, a ...interface{}) error {
	return &BadRequestError{msg:fmt.Sprintf(format, a...)}
}

