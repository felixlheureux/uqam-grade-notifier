package httperror

import (
	"fmt"
	"net/http"

	"github.com/felixlheureux/uqam-grade-notifier/pkg/domain"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// NewError generates function, so we can predefine a list of http errors without knowing the cause at compile time
func NewError(statusCode int, errorCode int, message string) func(error) *Error {
	return func(cause error) *Error {
		e := errors.WithStack(cause)

		return &Error{
			StatusCode:    statusCode,
			ErrorCode:     errorCode,
			OutputMessage: message,
			Cause:         e,
		}
	}
}

// FromEcho is similar to NewError but handles echo errors
func FromEcho(httpCode int, errorCode int, message string) func(*echo.HTTPError) *Error {
	return func(cause *echo.HTTPError) *Error {
		e := errors.WithStack(cause.Internal)

		return &Error{
			StatusCode:    httpCode,
			ErrorCode:     errorCode,
			OutputMessage: fmt.Sprintf("%s: %s", message, cause.Message),
			Cause:         e,
		}
	}
}

func FromDomain(err error) error {
	var result *domain.Error
	if dErr, ok := err.(*domain.Error); ok {
		result = dErr
	} else {
		result = domain.ErrUnexpected(err)
	}

	httpError := NewError(ErrStatusCode[result.Code], result.Code, result.Message)

	return httpError(result.Cause)
}

var (
	CoreEchoError                     = FromEcho(http.StatusInternalServerError, 0, "echo error")
	CoreUnknownError                  = NewError(http.StatusInternalServerError, 1, "unknown error")
	CoreRequestBindingFailed          = NewError(http.StatusBadRequest, 2, "failed to bind request body")
	CorePanic                         = NewError(http.StatusInternalServerError, 3, "panic")
	CoreDataUnmarshallFailed          = NewError(http.StatusBadRequest, 4, "failed to unmarshall data")
	CoreUnexpectedDataType            = NewError(http.StatusBadRequest, 5, "unexpected data type")
	CoreRequestFileFailed             = NewError(http.StatusBadRequest, 6, "failed to get file from request")
	CoreFileOpenFailed                = NewError(http.StatusInternalServerError, 7, "failed to open file")
	CoreRequestValidationFailed       = NewError(http.StatusBadRequest, 8, "failed to validate request")
	CoreRequestStringConversionFailed = NewError(http.StatusBadRequest, 9, "failed to convert string")
	CoreUnauthorized                  = NewError(http.StatusUnauthorized, 10, "unauthorized")
	CoreUnprocessableEntity           = NewError(http.StatusUnprocessableEntity, 11, "unprocessable entity")
)

var ErrStatusCode = map[int]int{
	domain.ErrUnexpected(nil).Code: http.StatusInternalServerError,
	domain.ErrNotFound(nil).Code:   http.StatusNotFound,

	domain.ErrUserGetFailed(nil).Code:    http.StatusInternalServerError,
	domain.ErrUserUpdateFailed(nil).Code: http.StatusInternalServerError,
}
