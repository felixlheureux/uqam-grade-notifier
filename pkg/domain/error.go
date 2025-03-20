package domain

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	// Core Errors (1000-1999)
	ErrUnexpected = NewError(1000, "An unexpected error occurred")
	ErrNotFound   = NewError(1001, "The requested resource does not exist")

	// Request Errors (2000-2999)
	CoreRequestBindingFailed    = NewError(2000, "Invalid request format")
	CoreRequestValidationFailed = NewError(2001, "Invalid request data")

	// Auth Errors (3000-3999)
	ErrTokenGenerationFailed = NewError(3000, "Failed to generate token")
	ErrTokenValidationFailed = NewError(3001, "Invalid or expired token")
	ErrTokenStorageFailed    = NewError(3002, "Failed to store token")
	ErrEmailValidationFailed = NewError(3003, "Invalid email")
	ErrEmailSendingFailed    = NewError(3004, "Failed to send email")

	// Course Errors (4000-4999)
	ErrCourseStorageFailed   = NewError(4000, "Failed to store courses")
	ErrGradeValidationFailed = NewError(4001, "Invalid grade")

	// User Errors (6000-6999)
	ErrUserGetFailed     = NewError(6000, "Failed to retrieve user")
	ErrUserFindOneFailed = NewError(6001, "User not found")
	ErrUserFindFailed    = NewError(6002, "Failed to search users")
	ErrUserCreateFailed  = NewError(6003, "Failed to create user")
	ErrUserUpdateFailed  = NewError(6004, "Failed to update user")
	ErrUserDeleteFailed  = NewError(6005, "Failed to delete user")
	ErrUserDestroyFailed = NewError(6006, "Failed to destroy user")
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Cause   error  `json:"cause"`
}

func NewError(code int, msg string) func(error) *Error {
	return func(cause error) *Error {
		e := errors.WithStack(cause)

		return &Error{
			Code:    code,
			Message: msg,
			Cause:   e,
		}
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("[%d] %s  %+v", e.Code, e.Message, e.Cause)
}

// IsNotFound vérifie si une erreur est de type NotFound
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}

	var domainErr *Error
	if errors.As(err, &domainErr) {
		return domainErr.Code == ErrNotFound(nil).Code
	}

	return false
}
