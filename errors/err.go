package errors

type ErrorType uint32

const (
	Ok ErrorType = iota
	WrongFile
	ResourceUnavailable
	InternalError
)

var Messages = map[ErrorType]string{
	WrongFile:           "Wrong File Format",
	ResourceUnavailable: "Resource Unavailable",
	InternalError:       "Internal server error",
}

type ServiceError struct {
	errorType ErrorType
}

func (se ServiceError) Error() string {
	return Messages[se.errorType]
}

func (se ServiceError) Type() ErrorType {
	return se.errorType
}

func NewServiceError(t ErrorType) ServiceError {
	return ServiceError{
		errorType: t,
	}
}
