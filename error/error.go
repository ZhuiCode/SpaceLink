package linkerr

import "fmt"

type LinkErr struct {
	code int64
	msg  string
}

var ErrList = []LinkErr{
	{0, "Path has been found"},
	{1, "Path has not been found"},
}

func NewLinkErr(code int64, msg string) LinkErr {
	return LinkErr{
		code: code,
		msg:  msg,
	}
}

func (e LinkErr) Error() string {
	return fmt.Sprintf("Error: [%d] %s", e.code, e.msg)
}

func (e LinkErr) GetCode() int64 {
	return e.code
}

func (e LinkErr) GetMsg() string {
	return e.msg
}
