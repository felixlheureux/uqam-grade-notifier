package httperror

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Response struct {
	ErrorCode int    `json:"error_code"`
	Error     string `json:"error"`
}

type Error struct {
	StatusCode int
	ErrorCode  int
	// message returned to the client
	OutputMessage string
	Cause         error
}

func (e Error) Error() string {
	return fmt.Sprintf("[%d] %s: %+v", e.ErrorCode, e.OutputMessage, e.Cause)
}

func (e Error) Is(other error) bool {
	err, ok := other.(*Error)

	if !ok {
		return false
	}

	return err.ErrorCode == e.ErrorCode
}

func NewErrorHandler(logger *log.Logger) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		var e *Error

		switch err := err.(type) {
		case *echo.HTTPError:
			e = CoreEchoError(err)
		case *Error:
			e = err
		default:
			e = CoreUnknownError(err)
		}

		if c.Response().Committed {
			logger.Printf("response already committed: %v", err)
			return
		}

		if c.Request().Method == http.MethodHead { // Issue https://github.com/labstack/echo/issues/608
			err = c.NoContent(e.StatusCode)
		} else {
			err = c.JSON(e.StatusCode, Response{e.ErrorCode, e.OutputMessage})
		}

		if err != nil {
			logger.Printf("unable to write error response: %v", err)
		}
	}
}
