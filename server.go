package godge

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"

	docker "github.com/fsouza/go-dockerclient"
)

type users map[string]string

// Server holds all the information related to a single instance of the judge. It's used to register Tasks and start the HTTP server.
type Server struct {
	address            string
	tasks              map[string]Task
	pendingSubmissions chan submissionRequest
	requestErrorChan   chan error
	dockerClient       *docker.Client
	users              users
	scoreboard         scoreboard
}

// NewServer creates a new instance of the judge. It takes the address that the
// judge will listen to and the address of the address daemon (e.g. unix:///var/run/docker.sock).
// NewServer returns an error if it fails to connect to the docker daemon.
func NewServer(address string, dockerAddress string) (*Server, error) {
	dc, err := docker.NewClient(dockerAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to docker daemon: %v", err)
	}
	if err := dc.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to docker daemon: %v", err)
	}
	return &Server{
		address:            address,
		tasks:              make(map[string]Task),
		pendingSubmissions: make(chan submissionRequest),
		dockerClient:       dc,
		users:              make(users),
		scoreboard:         make(scoreboard),
	}, nil
}

// RegisterTask registers a new task in the server.
func (s *Server) RegisterTask(t Task) {
	s.tasks[t.Name] = t
}

// handleSubmission is used to handle a received submission by executing the tests of the
// submission's task against this submission.
func (s *Server) handleSubmission(sub *Submission) error {
	t, ok := s.tasks[sub.TaskName]
	if !ok {
		return fmt.Errorf("task %v not found", sub.TaskName)
	}
	if err := t.execute(sub); err != nil {
		return fmt.Errorf("task %v failed: %v", sub.TaskName, err)
	}
	return nil
}

// A wrapper around the submission that's used for communication between
// the http handler and the server.
type submissionRequest struct {
	result     chan error
	submission *Submission
}

// Updates the scoreboard.
func (s *Server) reportResult(sub *Submission, err error) {
	log.Printf("%v submission for %v: %v", sub.Language, sub.TaskName, err)
	if err != nil {
		s.scoreboard.set(sub.Username, sub.TaskName, failedVerdict)
		return
	}
	s.scoreboard.set(sub.Username, sub.TaskName, passedVerdict)
}

// Executes the tests and report the result back to the http handler and the
// scoreboard.
func (s *Server) processSubmissions() {
	for sreq := range s.pendingSubmissions {
		err := s.handleSubmission(sreq.submission)
		sreq.result <- err
		s.reportResult(sreq.submission, err)
	}
}

// SubmissionResponse is the response returned back by the server in response
// to the submission request. It's exposed to be used by the command line client.
type SubmissionResponse struct {
	Passed bool   `json:"passed"`
	Error  string `json:"error"`
}

// The handler that handles submission requests.
func (s *Server) submitHTTPHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		httpJSONError(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}
	// Authenticate the user
	if username, password, ok := req.BasicAuth(); !ok || s.users[username] != password {
		httpJSONError(w, "Wrong username or password", http.StatusUnauthorized)
		return
	}

	var sub Submission
	err := json.NewDecoder(req.Body).Decode(&sub)
	if err != nil {
		httpJSONError(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
		return
	}
	sub.Executor.setDockerClient(s.dockerClient)

	// Send the submission for the server to run the tests.
	res := make(chan error)
	s.pendingSubmissions <- submissionRequest{
		result:     res,
		submission: &sub,
	}

	// Wait for the submission results and prepare the response.
	result := <-res

	resp := SubmissionResponse{
		Passed: true,
		Error:  "",
	}

	if result != nil {
		resp = SubmissionResponse{
			Passed: false,
			Error:  result.Error(),
		}
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		httpJSONError(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// RegisterRequest represents the registeration request. It's exposed to be used by the
// command line client.
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Handles registration requests.
func (s *Server) registerHTTPHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		httpJSONError(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	var rreq RegisterRequest
	err := json.NewDecoder(req.Body).Decode(&rreq)
	if err != nil {
		httpJSONError(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
		return
	}

	// Make sure that the username is not empty.
	if len(rreq.Username) == 0 {
		httpJSONError(w, fmt.Sprintf("Username cannot be empty"), http.StatusBadRequest)
		return
	}

	// Make sure that the username is unique.
	if _, ok := s.users[rreq.Username]; ok {
		httpJSONError(w, fmt.Sprintf("Username %v is already registered", rreq.Username), http.StatusBadRequest)
		return
	}

	s.users[rreq.Username] = rreq.Password
	log.Printf("User %v registered", rreq.Username)

	w.WriteHeader(http.StatusCreated)
}

// Handles tasks queries.
func (s *Server) tasksHTTPHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		httpJSONError(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	var ts []Task

	for _, t := range s.tasks {
		ts = append(ts, t)
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(ts); err != nil {
		httpJSONError(w, "Failed to encode tasks", http.StatusInternalServerError)
		return
	}
}

// Handles scoreboard requests.
func (s *Server) scoreboardHTTPHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		httpJSONError(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Add("Content-Type", "text/html")

	var ts []string
	for k := range s.tasks {
		ts = append(ts, k)
	}
	sort.Strings(ts)

	var us []string
	for k := range s.users {
		us = append(us, k)
	}
	sort.Strings(us)

	scoreboard := s.scoreboard.toScoreboard(us, ts)

	scoreboardTmpl.Execute(w, map[string]interface{}{
		"Scoreboard": scoreboard,
	})
}

// Start starts the http server and the goroutine responsible for processing
// the submissions.
func (s *Server) Start() error {
	go s.processSubmissions()
	mux := http.NewServeMux()
	mux.HandleFunc("/submit", s.submitHTTPHandler)
	mux.HandleFunc("/register", s.registerHTTPHandler)
	mux.HandleFunc("/tasks", s.tasksHTTPHandler)
	mux.HandleFunc("/scoreboard", s.scoreboardHTTPHandler)
	return http.ListenAndServe(s.address, mux)
}
