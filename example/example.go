package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/MohamedBassem/godge"
)

func runAndCompareOutput(sub *godge.Submission, args []string, want string) error {
	if err := sub.Executor.Execute(args); err != nil {
		return err
	}
	defer sub.Executor.Stop()
	<-sub.Executor.DieEvent()
	got, err := sub.Executor.Stdout()
	if err != nil {
		return err
	}
	if got != want {
		return fmt.Errorf("want: %v, got: %v", want, got)
	}
	return nil
}

var tasks = []godge.Task{
	{
		Name: "HelloWorld",
		Desc: "Your program should print 'Hello World!' to stdout.",
		Tests: []godge.Test{
			{
				Name: "PrintsHelloWorld",
				Func: func(sub *godge.Submission) error {
					return runAndCompareOutput(sub, []string{}, "Hello World!")
				},
			},
		},
	},
	{
		Name: "Flags",
		Desc: "The binary should accept the '--name' flag and prints 'Hello $name!' to stdout.",
		Tests: []godge.Test{
			{
				Name: "PrintsHelloJudge",
				Func: func(sub *godge.Submission) error {
					return runAndCompareOutput(sub, []string{"--name", "Judge"}, "Hello Judge!")
				},
			},
			{
				Name: "PrintsHelloUser",
				Func: func(sub *godge.Submission) error {
					return runAndCompareOutput(sub, []string{"--name", sub.Username}, fmt.Sprintf("Hello %v!", sub.Username))
				},
			},
		},
	},
}

func main() {
	rand.Seed(time.Now().UnixNano())
	server, err := godge.NewServer(":8080", "unix:///var/run/docker.sock", "tmp.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	for _, t := range tasks {
		server.RegisterTask(t)
	}
	log.Fatal(server.Start())
}
