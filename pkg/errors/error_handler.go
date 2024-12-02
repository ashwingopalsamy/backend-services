package errors

import "fmt"

type CustomError struct {
	StatusCode int
	Message    string
	Details    map[string]string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("status_code: %d, message: %s", e.StatusCode, e.Message)
}

func New(statusCode int, message string, details map[string]string) *CustomError {
	return &CustomError{
		StatusCode: statusCode,
		Message:    message,
		Details:    details,
	}
}

func BadRequest(message string, details map[string]string) *CustomError {
	return New(400, message, details)
}

//func Unauthorized(message string) *CustomError {
//	return New(401, message, nil)
//}
//
//func Conflict(message string) *CustomError {
//	return New(409, message, nil)
//}
//
//func UnprocessableEntity(message string) *CustomError {
//	return New(422, message, nil)
//}
//
//func Locked(message string) *CustomError {
//	return New(423, message, nil)
//}
//
//func InternalServerError(message string) *CustomError {
//	return New(500, message, nil)
//}
