package main

import (
	"fmt"
	"log"
	"time"

	"github.com/MohamedBassem/godge"
)

var tasks = []godge.Task{
	{
		Name: "HelloWorld",
		Desc: "Your program should print 'Hello World!' to stdout.",
		Tests: []godge.Test{
			{
				Name: "PrintsHelloWorld",
				Func: func(sub *godge.Submission) error {
					if err := sub.Executor.Execute([]string{}); err != nil {
						return err
					}
					defer sub.Executor.Stop()
					time.Sleep(1 * time.Second)
					want := "Hello World!"
					got, err := sub.Executor.Stdout()
					if err != nil {
						return err
					}
					if got != want {
						return fmt.Errorf("want: %v, got: %v", want, got)
					}
					return nil
				},
			},
		},
	},
	{
		Name: "WebServer",
		Tests: []godge.Test{
			{
				Name: "RespondsWith200",
				Func: func(sub *godge.Submission) error {
					return nil
				},
			},
		},
	},
}

func main() {
	server, err := godge.NewServer(":8080", "unix:///var/run/docker.sock")
	if err != nil {
		log.Fatal(err)
	}
	for _, t := range tasks {
		server.RegisterTask(t)
	}
	log.Fatal(server.Start())
}
