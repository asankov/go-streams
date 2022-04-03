package stream

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMultiErrorEmpty(t *testing.T) {
	testCases := []struct {
		name   string
		errors []error
	}{
		{
			name:   "nil array",
			errors: nil,
		},
		{
			name:   "empty array",
			errors: []error{},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			e := &MultiError{errors: testCase.errors}

			require.Empty(t, e.Error())
		})
	}
}

func TestMultiErrorOneError(t *testing.T) {
	msg := "only one error"
	e := &MultiError{errors: []error{fmt.Errorf(msg)}}

	require.Equal(t, msg, e.Error())
}

func TestMultiErrorMultipleErrors(t *testing.T) {
	messages := []string{
		"first msg",
		"second msg",
		"third msg",
	}

	errors := make([]error, 0, len(messages))
	expectedMsg := "multiple errors: "
	for _, msg := range messages {
		errors = append(errors, fmt.Errorf(msg))
		expectedMsg += msg + ", "
	}

	e := &MultiError{errors: errors}

	require.Equal(t, expectedMsg, e.Error())
}
