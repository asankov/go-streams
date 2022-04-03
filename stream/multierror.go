package stream

import "strings"

type MultiError struct {
	errors []error
}

// Error implements the error interface.
func (e *MultiError) Error() string {
	if len(e.errors) == 0 {
		return ""
	}
	if len(e.errors) == 1 {
		return e.errors[0].Error()
	}
	var sb strings.Builder
	sb.WriteString("multiple errors: ")
	for _, err := range e.errors {
		_, _ = sb.WriteString(err.Error())
		_, _ = sb.WriteString(", ")
	}
	return sb.String()
}
