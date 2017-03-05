package godge

import (
	"encoding/json"
	"fmt"
)

type Submission struct {
	Language string   `json:"language"`
	TaskName string   `json:"taskName"`
	Executor Executor `json:"submission"`
}

func (s *Submission) UnmarshalJSON(d []byte) error {
	metadata := struct {
		Language   string          `json:"language"`
		TaskName   string          `json:"taskName"`
		Submission json.RawMessage `json:"submission"`
	}{}

	err := json.Unmarshal(d, &metadata)
	if err != nil {
		return fmt.Errorf("failed to unmarshal metadata: %v", err)
	}

	s.Language = metadata.Language
	s.TaskName = metadata.TaskName

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
