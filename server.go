package godge

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	docker "github.com/fsouza/go-dockerclient"
)

type users map[string]string

type Server struct {
	address            string
	tasks              map[string]Task
	pendingSubmissions chan *Submission
	requestErrorChan   chan error
	dockerClient       *docker.Client
	users              users
}

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
		pendingSubmissions: make(chan *Submission, 20),
		dockerClient:       dc,
		users:              make(users),
	}, nil
}

func (s *Server) RegisterTask(t Task) {
	s.tasks[t.Name] = t
}

func (s *Server) handleSubmission(sub *Submission) error {
	t, ok := s.tasks[sub.TaskName]
	if !ok {
		return fmt.Errorf("task %v not found", sub.TaskName)
	}
	if err := t.Execute(sub); err != nil {
		return fmt.Errorf("task %v failed: %v", sub.TaskName, err)
	}
	return nil
}

func (s *Server) reportResult(sub *Submission, err error) {
	// TODO(mbassem): Report actual result to the scoreboard
	log.Printf("%v submission for %v: %v", sub.Language, sub.TaskName, err)
}

func (s *Server) processSubmissions() {
	for sub := range s.pendingSubmissions {
		s.reportResult(sub, s.handleSubmission(sub))
	}
}

func (s *Server) submitHTTPHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		httpJsonError(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}
	if username, password, ok := req.BasicAuth(); !ok || s.users[username] != password {
		httpJsonError(w, "Wrong username or password", http.StatusUnauthorized)
		return
	}
	var sreq Submission
	err := json.NewDecoder(req.Body).Decode(&sreq)
	if err != nil {
		httpJsonError(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
		return
	}
	sreq.Executor.setDockerClient(s.dockerClient)

	s.pendingSubmissions <- &sreq
	w.WriteHeader(http.StatusNoContent)
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *Server) registerHTTPHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		httpJsonError(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	var rreq RegisterRequest
	err := json.NewDecoder(req.Body).Decode(&rreq)
	if err != nil {
		httpJsonError(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
		return
	}

	if len(rreq.Username) == 0 {
		httpJsonError(w, fmt.Sprintf("Username cannot be empty"), http.StatusBadRequest)
		return
	}

	if _, ok := s.users[rreq.Username]; ok {
		httpJsonError(w, fmt.Sprintf("Username %v is already registered", rreq.Username), http.StatusBadRequest)
		return
	}

	s.users[rreq.Username] = rreq.Password
	log.Printf("User %v registered", rreq.Username)

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) tasksHTTPHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		httpJsonError(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	var ts []Task

	for _, t := range s.tasks {
		ts = append(ts, t)
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(ts); err != nil {
		httpJsonError(w, "Failed to encode tasks", http.StatusInternalServerError)
		return
	}
}

func (s *Server) Start() error {
	go s.processSubmissions()
	mux := http.NewServeMux()
	mux.HandleFunc("/submit", s.submitHTTPHandler)
	mux.HandleFunc("/register", s.registerHTTPHandler)
	mux.HandleFunc("/tasks", s.tasksHTTPHandler)
	return http.ListenAndServe(s.address, mux)
}
