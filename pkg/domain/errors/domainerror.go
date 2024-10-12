package errors

type ErrorCode string

type Domain interface {
	error
	Code() ErrorCode
}
