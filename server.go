package godge

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	address            string
	tasks              map[string]Task
	pendingSubmissions chan *Submission
	requestErrorChan   chan error
}

func NewServer(address string) *Server {
	return &Server{
		address:            address,
		tasks:              make(map[string]Task),
		pendingSubmissions: make(chan *Submission),
	}
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

func (s *Server) Start() error {
	go s.processSubmissions()
	mux := http.NewServeMux()
	mux.HandleFunc("/submit", submitHandler)
	return http.ListenAndServe(s.address, mux)
}
