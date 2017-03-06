package godge

import "fmt"

type Test struct {
	Name string
	Func func(*Submission) error
}

type Task struct {
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Tests []Test `json:"-"`
}

func (t *Task) Execute(s *Submission) error {
	var errs Errors
	for _, test := range t.Tests {
		if err := test.Func(s); err != nil {
			errs = append(errs, fmt.Errorf("test '%v' failed: %v", test.Name, err))
		}
	}
	return errs.ErrorOrNil()
}
