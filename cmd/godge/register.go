package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/MohamedBassem/godge"
	"github.com/google/subcommands"
)

type registerCmd struct {
	username string
	password string
}

func (*registerCmd) Name() string     { return "register" }
func (*registerCmd) Synopsis() string { return "Registers a new user." }
func (*registerCmd) Usage() string {
	return `register -username <username> -password <password>:
  Registers a new user.
`
}

func (s *registerCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&s.username, "username", "", "The username you want to register with")
	f.StringVar(&s.password, "password", "", "The password you want to register with")
}

func (s *registerCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
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

	req := godge.RegisterRequest{
		Username: s.username,
		Password: s.password,
	}

	reqj, err := json.Marshal(&req)
	if err != nil {
		log.Printf("failed to marshal request json: %v", err)
		return subcommands.ExitFailure
	}

	resp, err := http.Post(fmt.Sprintf("%v/register", *serverAddress), "application/json", bytes.NewReader(reqj))
	if err != nil {
		log.Printf("failed to submit request: %v", err)
		return subcommands.ExitFailure
	}
	defer resp.Body.Close()
	if err := checkResponseError(resp); err != nil {
		log.Printf("Registeration failed: %v", err)
		return subcommands.ExitFailure
	}

	log.Println("Registered ..")
	return subcommands.ExitSuccess
}
