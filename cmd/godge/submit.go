package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/MohamedBassem/godge"
	"github.com/google/subcommands"
)

type submitCmd struct {
	language string
	taskName string
	username string
	password string
}

func (*submitCmd) Name() string     { return "submit" }
func (*submitCmd) Synopsis() string { return "Submits solution to the server." }
func (*submitCmd) Usage() string {
	return `submit -languge <language> -task <taskName> -username <username> -password <password>:
  Submits solution to the server.
`
}

func (s *submitCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&s.language, "language", "", "The language of the submission")
	f.StringVar(&s.taskName, "task", "", "The task of the submission")
	f.StringVar(&s.username, "username", "", "Your username")
	f.StringVar(&s.password, "password", "", "Your password")
}

func (s *submitCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if s.language == "" {
		log.Println("Language must be specified")
		return subcommands.ExitUsageError
	}
	if s.taskName == "" {
		log.Println("Task must be specified")
		return subcommands.ExitUsageError
	}
	if s.username == "" {
		log.Println("Username must be specified")
		return subcommands.ExitUsageError
	}
	if s.password == "" {
		log.Println("Password must be specified")
		return subcommands.ExitUsageError
	}
	if *serverAddress == "" {
		log.Fatal("Server Address must be specified")
	}

	switch s.language {
	case "go":
		if err := s.submit(s.goSubmission); err != nil {
			log.Println(err)
			return subcommands.ExitFailure
		}
	default:
		log.Printf("Unsupported language: %v", s.language)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

func (s *submitCmd) submit(executor func() (godge.Executor, error)) error {

	exec, err := executor()
	if err != nil {
		return fmt.Errorf("failed to create executor: %v", err)
	}

	sub := &godge.Submission{
		Language: s.language,
		TaskName: s.taskName,
		Username: s.username,
		Executor: exec,
	}

	reqj, err := json.Marshal(&sub)
	if err != nil {
		return fmt.Errorf("failed to marshal request json: %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%v/submit", *serverAddress), bytes.NewReader(reqj))
	if err != nil {
		return fmt.Errorf("failed to marshal request json: %v", err)
	}
	req.SetBasicAuth(s.username, s.password)

	client := &http.Client{
		Timeout: 0,
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to submit request: %v", err)
	}
	defer resp.Body.Close()
	if err := checkResponseError(resp); err != nil {
		return fmt.Errorf("submission failed: %v", err)
	}

	var result godge.SubmissionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	if result.Passed {
		log.Println("You submission passed!")
	} else {
		log.Printf("You submission failed: %v", result.Error)
	}

	return nil
}

func (s *submitCmd) goSubmission() (godge.Executor, error) {

	// Zip the current dir
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working dir: %v", err)
	}
	log.Printf("Will submit %v", currentDir)
	b, err := zipCurrentDir()
	if err != nil {
		return nil, fmt.Errorf("failed to zip current dir: %v", err)
	}
	log.Printf("Done zipping %v", currentDir)

	return &godge.GoExecutor{
		PackageArchive: b,
	}, nil
}
