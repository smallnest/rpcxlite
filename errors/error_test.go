package errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestNewMultiError(t *testing.T) {
	var errorSet []error
	errorSet = append(errorSet, errors.New("invalid"))
	errorSet = append(errorSet, errors.New("fatal"))

	multiError := NewMultiError(errorSet)
	if fmt.Sprintf("%v", errorSet) != multiError.Error() {
		t.Fatal("Test NewMultiError()")
	}
}

func TestMultiError_Append(t *testing.T) {
	multiErrors := MultiError{}
	multiErrors.Errors = append(multiErrors.Errors, errors.New("invalid"))
	multiErrors.Errors = append(multiErrors.Errors, errors.New("fatal"))

	if len(multiErrors.Errors) != 2 {
		t.Fatal("Test Append()")
	}
}

func TestMultiError_Error(t *testing.T) {
	multiErrors := MultiError{}
	multiErrors.Errors = append(multiErrors.Errors, errors.New("invalid"))
	multiErrors.Errors = append(multiErrors.Errors, errors.New("fatal"))

	if "[invalid fatal]" != multiErrors.Error() {
		t.Fatal("Test Error()")
	}
}
