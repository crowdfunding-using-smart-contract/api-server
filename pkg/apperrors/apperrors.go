package apperrors

type appError struct {
	Code    int
	Message string
}

type Error interface {
	Status() int
	Error() string
}

func New(code int, message string) Error {
	return &appError{
		Code:    code,
		Message: message,
	}
}

func (e *appError) Status() int {
	return e.Code
}

func (e *appError) Error() string {
	return e.Message
}
