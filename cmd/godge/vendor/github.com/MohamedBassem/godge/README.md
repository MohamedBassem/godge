# Godge

Godge is a self-hosted online judge for meetups and workshops. It adds a competitive environment to the workshop/meetup.

## Why I built Godge?

Read the blog post : [http://blog.mbassem.com/2017/03/12/godge/](http://blog.mbassem.com/2017/03/12/godge/).

## What's the difference between godge and other online judges?

1- The obvious difference is that godge is self hosted. It's not intended to be running for public internet audience. It targets meetup/workshop attendees.

2- With Godge you can test every submission the way you want. Other online judges like (spoj and codeforces) only test that the user's submission printed the correct output. When you are teaching Go for instance, you want to test that attendees can create command line flags, that they can start a webserver and respond 200 to your requests, etc. Godge allows you do so.

3- Other online judges only allow you to submit a single file. Godge allows you test the project as a whole, with its subpackages, config files, assets, etc.

## How to use it?

### As a meetup/workshop Host

1- Install docker (> 1.22).

2- Implement the tasks that you want your attendees to solve using the Godge package. An example for two tasks that tests that attendees can print output to stdout and that they can use flags correctly.

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
```

3- Share with your attendees the address of the server. You can host it on the local network or on a public server.

### As an Attendee

1- Download the Godge command line client.

```
$ go get -u github.com/MohamedBassem/godge/cmd/godge
```

2- Register a new account on the server.

```
$ godge --address <addr> register --username <username> --password <password>
```

3- List the available tasks.

```
$ godge --address <addr> tasks
HelloWorld: Your program should print 'Hello World!' to stdout.
=============================
Flags: The binary should accept the '--name' flag and prints 'Hello $name!' to stdout.
=============================
```

4- Work on one of the tasks and then submit it.

```
$ godge --address <addr> submit --task <task> --language <lang> --username <username> --password <password>
2017/03/12 19:04:58 Will submit /private/tmp/tmp
2017/03/12 19:04:58 Done zipping /private/tmp/tmp
2017/03/12 19:05:00 You submission passed!
```

5- Check the scoreboard at `http://<addr>/scoreboard`.

## A Live Demo

[![Demo Image](https://raw.githubusercontent.com/MohamedBassem/godge/master/demo_image.png)](https://www.youtube.com/watch?v=S0OLOiujxJk)

## How It Works

Submissions run in a separate container. The container is determined based on the language. Godge offers
an abstract API to interact with the container (start, stop, fetch stdout, ..).

### Go

The command line client, zips the whole "main" package (and its subpackages) and sends it to the server. The server
then uses the image `golang:1.8` to build and run the package.

## Known Issues / Future Work

1- Currently `Go` is the only supported language.

~~2- All the data (scoreboard, tasks and users) are stored in the server's memory. All the
data will be lost if the server is restarted. Although it's not a problem for meetups or
workshop as they are short by nature, it would be nice to persist this info in a `sqlite`
database for example.~~

3- The `Execute` function should allow opening ports in the container to be able to test
web servers for example.

4- Currently a single goroutine executes the submissions sequentially. It would be nice
to run multiple submissions in parallel. [Easy Fix]

##Contribution

Your contributions and ideas are welcomed through issues and pull requests.

##License

Copyright (c) 2017, Mohamed Bassem. (MIT License)

See LICENSE for more info.
