package godge

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	docker "github.com/fsouza/go-dockerclient"
)

type Server struct {
	address            string
	tasks              map[string]Task
	pendingSubmissions chan *Submission
	requestErrorChan   chan error
	dockerClient       *docker.Client
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

	// TODO(mbassem): Authenticate the request

	defer req.Body.Close()

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

func (s *Server) Start() error {
	go s.processSubmissions()
	mux := http.NewServeMux()
	mux.HandleFunc("/submit", s.submitHTTPHandler)
	return http.ListenAndServe(s.address, mux)
}
