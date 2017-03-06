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
		Tests: []godge.Test{
			{
				Name: "PrintsHelloWorld",
				Func: func(sub *godge.Submission) error {
					err := sub.Executor.Execute([]string{})
					if err != nil {
						return err
					}
					defer sub.Executor.Stop()
					time.Sleep(1 * time.Second)
					fmt.Println(sub.Executor.Stdout())
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
