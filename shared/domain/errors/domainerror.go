package errors

type ErrorCode string
type ErrorType string

const (
	BusinessErrorType ErrorType = "Business"
	SystemErrorType   ErrorType = "System"
)

type Domain interface {
	error
	Code() ErrorCode
	Type() ErrorType
}
