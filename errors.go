package godge

import "fmt"

// Errors represents a collection of errors.
type Errors []error

// Error returns the error string.
func (e Errors) Error() string {
	return fmt.Sprintf("%v", []error(e))
}

// ErrorOrNil returns nill if there is no error, otherwise returns the error.
func (e Errors) ErrorOrNil() error {
	if len(e) != 0 {
		return e
	}
	return nil
}
