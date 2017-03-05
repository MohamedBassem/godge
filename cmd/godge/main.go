package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/MohamedBassem/godge"
)

var (
	language      = flag.String("language", "", "The language of the submission")
	taskName      = flag.String("task", "", "The task of the submission")
	serverAddress = flag.String("address", "", "The address of the server")
)

func goSubmission() error {

	// Zip the current dir
	b, err := zipCurrentDir()
	if err != nil {
		return fmt.Errorf("failed to zip current dir: %v", err)
	}

	req := godge.Submission{
		Language: "go",
		TaskName: *taskName,
		Executor: &godge.GoExecutor{
			PackageArchive: b,
		},
	}

	reqj, err := json.Marshal(&req)
	if err != nil {
		return fmt.Errorf("failed to marshal request json: %v", err)
	}

	fmt.Println(string(reqj))

	resp, err := http.Post(fmt.Sprintf("%v/submit", *serverAddress), "application/json", bytes.NewReader(reqj))
	if err != nil {
		return fmt.Errorf("failed to submit request: %v", err)
	}
	defer resp.Body.Close()

	return nil
}

func main() {
	flag.Parse()

	if *language == "" {
		log.Fatal("Language must be specified")
	}
	if *taskName == "" {
		log.Fatal("Task must be specified")
	}
	if *serverAddress == "" {
		log.Fatal("Server Address must be specified")
	}

	switch *language {
	case "go":
		if err := goSubmission(); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("Unsupported language: %v", *language)
	}

	log.Println("Submitted ..")
}
