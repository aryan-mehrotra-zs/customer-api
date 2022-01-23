package errors

type BindError struct{}

func (e BindError) Error() string {
	return "bind error"
}
