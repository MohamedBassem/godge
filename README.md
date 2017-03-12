# Godge

Godge is an self-hosted online judge for meetups and workshops. It adds a competitive environment to the workshop/meetup.

## What's the difference between godge and other online judges?

1- The obvious difference is that godge is self hosted. It's not intended to be running for public internet audience. It targets meetup/workshop attendees.
2- With godge you can test every submission the way you want. Other online judges like (spoj and codeforces) only test that the user's submission printed the correct output. When you are teaching Go for instance, you want to test that attendees can create command line flags, that they can start a webserver and respond 200 to your requests, etc. Godge allows you do so.
3- Other online judges only allow you to submit a single file. Godge allows you test the project as a whole, with its subpackages, config files, assets, etc.

## How to use it?

### As a meetup/workshop Hosts

1- Install docker.
2- Implement the tasks that you want your attendees to solve using the godge package. An example for two tasks that tests that attendees can print output to stdout and that they can use flags correctly.

```go
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
	time.Sleep(1 * time.Second)
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
	server, err := godge.NewServer(":8080", "unix:///var/run/docker.sock")
	if err != nil {
		log.Fatal(err)
	}
	for _, t := range tasks {
		server.RegisterTask(t)
	}
	log.Fatal(server.Start())
}
```

3- Share with your attendees the address of the server. You can host it on the local network or on a public server.

### As an Attendee
