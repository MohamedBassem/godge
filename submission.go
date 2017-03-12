package godge

import (
	"encoding/json"
	"fmt"
)

// Submission is the input of the user defined task tests.
type Submission struct {
	// The language of the submission.
	Language string `json:"language"`
	// The task this submission is sent to.
	TaskName string `json:"taskName"`
	// The username of the submitter.
	Username string `json:"username"`
	// The executor interface to deal with the submission.
	Executor Executor `json:"submission"`
}

// UnmarshalJSON is a custom JSON unmarshaller. It's used mainly to create
// a new executor instance based on the language field of the submission.
func (s *Submission) UnmarshalJSON(d []byte) error {
	metadata := struct {
		Language   string          `json:"language"`
		TaskName   string          `json:"taskName"`
		Username   string          `json:"username"`
		Submission json.RawMessage `json:"submission"`
	}{}

	err := json.Unmarshal(d, &metadata)
	if err != nil {
		return fmt.Errorf("failed to unmarshal metadata: %v", err)
	}

	s.Language = metadata.Language
	s.TaskName = metadata.TaskName
	s.Username = metadata.Username

	switch s.Language {
	case "go":
		var e GoExecutor
		err := json.Unmarshal(metadata.Submission, &e)
		if err != nil {
			return fmt.Errorf("failed to unmarshal language specific json: %v", err)
		}
		s.Executor = &e
	default:
		return fmt.Errorf("unsupported language %v", s.Language)
	}

	return nil
}
