package main

import (
	"flag"
	"fmt"
	"os"

	"tools.icpc.global/contest/commands"
)

// Current use:
//   contest -url <baseURL> -u <user> -p <password> contests
//   contest -url <baseURL> -c contestId -u <user> -p <password> problems
//   contest -url <baseURL> -c contestId -u <user> -p <password> clar <text>
//   contest -url <baseURL> -c contestId -u <user> -p <password> submit <text>
// -insecure is optional, the rest isn't. Problem id is hardcoded to "checks" for now

var (
	baseURL   string
	contestId string
	user      string
	password  string

	insecure bool
)

func init() {
	flag.StringVar(&baseURL, "url", "", "the base 'Contest URL'")
	flag.StringVar(&contestId, "c", "", "the 'contest id'")
	flag.StringVar(&user, "u", "", "the 'user'")
	flag.StringVar(&password, "p", "", "the 'password'")
	flag.BoolVar(&insecure, "insecure", false, "whether the connection is secure")
}

func main() {
	err := commands.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
