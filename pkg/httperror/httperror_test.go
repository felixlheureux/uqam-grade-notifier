package httperror

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsErrors(t *testing.T) {
	err1 := CoreUnknownError(errors.New("test my error"))
	assert.True(t, errors.Is(err1, CoreUnknownError(nil)))
}
