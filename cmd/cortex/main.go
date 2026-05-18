package main

import (
	"fmt"
	"os"

	"github.com/AdeshDeshmukh/cortex/cmd/cortex/commands"
)

var version = "dev"

func main() {
	commands.Version = version

	if err := commands.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
