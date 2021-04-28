package main

import (
	"os"

	"tools.icpc.global/contest/commands"
)

func main() {
	err := commands.Execute()
	if err != nil {
		os.Exit(1)
	}
}
