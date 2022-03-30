package main

import (
	"os"

	"github.com/icpctools/cli/commands"
)

func main() {
	err := commands.Execute()
	if err != nil {
		os.Exit(1)
	}
}
