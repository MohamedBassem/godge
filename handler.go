package godge

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func httpJsonError(w http.ResponseWriter, msg string, code int) {
	e := struct {
		Error string `json:"error"`
	}{
		Error: msg,
	}

	b, _ := json.Marshal(e)
	http.Error(w, string(b), code)
}

func submitHandler(ch chan<- *Submission) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
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

		ch <- &sreq
		w.WriteHeader(http.StatusNoContent)
	}
}
