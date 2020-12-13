package collection

import "fmt"

// NotAllowedf TODO
func NotAllowedf(format string, a ...interface{}) error {
	return notAllowedError{fmt.Sprintf(format, a...)}
}

type notAllowedError struct {
	msg string
}

func (e notAllowedError) Type() string {
	return "ClientError"
}

func (e notAllowedError) Error() string {
	return e.msg
}
