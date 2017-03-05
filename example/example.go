package main

import (
	"log"

	"github.com/MohamedBassem/godge"
)

var tasks = []godge.Task{
	{
		Name: "HelloWorld",
		Tests: []godge.Test{
			{
				Name: "PrintsHelloWorld",
				Func: func(sub *godge.Submission) error {
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
