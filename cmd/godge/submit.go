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
}

func (*submitCmd) Name() string     { return "submit" }
func (*submitCmd) Synopsis() string { return "Submits solution to the server." }
func (*submitCmd) Usage() string {
	return `submit -languge <language> -task <taskName>:
  Submits solution to the server.
`
}

func (s *submitCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&s.language, "language", "", "The language of the submission")
	f.StringVar(&s.taskName, "task", "", "The task of the submission")
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
	if *serverAddress == "" {
		log.Fatal("Server Address must be specified")
	}
	switch s.language {
	case "go":
		if err := s.goSubmission(); err != nil {
			log.Println(err)
			return subcommands.ExitFailure
		}
	default:
		log.Printf("Unsupported language: %v", s.language)
		return subcommands.ExitFailure
	}

	log.Println("Submitted ..")
	return subcommands.ExitSuccess
}

func (s *submitCmd) goSubmission() error {

	// Zip the current dir
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working dir: %v", err)
	}
	log.Printf("Will submit %v", currentDir)
	b, err := zipCurrentDir()
	if err != nil {
		return fmt.Errorf("failed to zip current dir: %v", err)
	}
	log.Printf("Done zipping %v", currentDir)

	req := godge.Submission{
		Language: "go",
		TaskName: s.taskName,
		Executor: &godge.GoExecutor{
			PackageArchive: b,
		},
	}

	reqj, err := json.Marshal(&req)
	if err != nil {
		return fmt.Errorf("failed to marshal request json: %v", err)
	}

	resp, err := http.Post(fmt.Sprintf("%v/submit", *serverAddress), "application/json", bytes.NewReader(reqj))
	if err != nil {
		return fmt.Errorf("failed to submit request: %v", err)
	}
	if err := checkResponseError(resp); err != nil {
		return fmt.Errorf("submission failed: %v", err)
	}
	defer resp.Body.Close()

	return nil
}
