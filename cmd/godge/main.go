package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"
)

var serverAddress = flag.String("address", "", "The address of the server")

func main() {
	subcommands.ImportantFlag("address")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&submitCmd{}, "")
	subcommands.Register(&registerCmd{}, "")
	subcommands.Register(&tasksCmd{}, "")
	flag.Parse()

	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
