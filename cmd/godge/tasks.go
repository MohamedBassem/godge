package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/MohamedBassem/godge"
	"github.com/google/subcommands"
)

type tasksCmd struct {
}

func (*tasksCmd) Name() string     { return "tasks" }
func (*tasksCmd) Synopsis() string { return "Prints the tasks and their description." }
func (*tasksCmd) Usage() string {
	return `tasks:
  Prints the tasks and their description.
`
}

func (s *tasksCmd) SetFlags(f *flag.FlagSet) {
}

func (s *tasksCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if *serverAddress == "" {
		log.Fatal("Server Address must be specified")
	}

	resp, err := http.Get(fmt.Sprintf("%v/tasks", *serverAddress))
	if err != nil {
		log.Printf("Failed to fetch tasks: %v", err)
		return subcommands.ExitFailure
	}
	defer resp.Body.Close()
	if err := checkResponseError(resp); err != nil {
		log.Printf("Fetching tasks failed: %v", err)
		return subcommands.ExitFailure
	}

	var ts []godge.Task

	if err := json.NewDecoder(resp.Body).Decode(&ts); err != nil {
		log.Printf("Decoding response failed: %v", err)
		return subcommands.ExitFailure
	}

	for _, t := range ts {
		fmt.Printf("%v: %v\n", t.Name, t.Desc)
		fmt.Println("=============================")
	}

	return subcommands.ExitSuccess
}
