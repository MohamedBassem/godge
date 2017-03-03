package godge

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func httpJsonError(w http.ResponseWriter, msg string, code int) {
	e := struct {
		Error string
	}{
		Error: msg,
	}

	b, _ := json.Marshal(e)
	http.Error(w, string(b), code)
}

type Submission struct {
	Language string
	TaskName string
	Executor Executor
}

func (s *Submission) UnmarshalJSON(d []byte) error {
	metadata := struct {
		Language         string
		TaskName         string
		LanguageSpecific json.RawMessage
	}{}

	err := json.Unmarshal(d, &metadata)
	if err != nil {
		return fmt.Errorf("failed to unmarshal metadata: %v", err)
	}

	s.Language = metadata.Language
	s.TaskName = metadata.TaskName

	switch s.Language {
	case "go":
		var e goExecutor
		err := json.Unmarshal(metadata.LanguageSpecific, &e)
		if err != nil {
			return fmt.Errorf("failed to unmarshal language specific json: %v", err)
		}
		s.Executor = &e
	default:
		return fmt.Errorf("unsupported language %v", s.Language)
	}

	return nil
}

func submitHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		httpJsonError(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	// TODO(mbassem): Authenticate the request

	defer req.Body.Close()

	var sreq Submission
	err := json.NewDecoder(req.Body).Decode(&sreq)
	if err != nil {
		httpJsonError(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
		return
	}

}
