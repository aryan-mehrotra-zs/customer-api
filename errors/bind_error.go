package errors

import "fmt"

type BindError struct{}

func (e BindError) Error() string {
	return fmt.Sprint("bind error")
}
