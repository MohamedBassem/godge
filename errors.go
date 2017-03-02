package godge

import "fmt"

type Errors []error

func (e Errors) Error() string {
	return fmt.Sprintf("%v", []error(e))
}

func (e Errors) ErrorOrNil() error {
	if len(e) != 0 {
		return e
	}
	return nil
}
