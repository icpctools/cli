package main

import (
	"fmt"
	"os"

	"tools.icpc.global/contest/commands"
)

func main() {
	err := commands.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
