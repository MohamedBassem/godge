package godge

import "fmt"

// Test defines on of the tests of a certain task.
type Test struct {
	// The name of the test, which will be returned to the user when the test
	// fails.
	Name string
	// The actuall test. It takes a submission as an input (along with its excutor)
	// and should return a descriptive error when the submission don't pass the test.
	Func func(*Submission) error
}

// Task defines a group of related tests. The user needs to pass all the tests to pass
// the task and get its point on the scoreboard.
type Task struct {
	// The name of the task that the user will use to submit their submission.
	Name string `json:"name"`
	// A description of what's required in order to pass the task.
	Desc string `json:"desc"`
	// A group of tests that a submission needs to pass in order to pass the task.
	Tests []Test `json:"-"`
}

// Execute runs the submission against all the tests. The error returned is the error
// retured by all the tests.
func (t *Task) execute(s *Submission) error {
	var errs Errors
	for _, test := range t.Tests {
		if err := test.Func(s); err != nil {
			errs = append(errs, fmt.Errorf("test '%v' failed: %v", test.Name, err))
		}
	}
	return errs.ErrorOrNil()
}
