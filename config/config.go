package config

import (
	"flag"
)

var (
	BaseURL   string
	ContestId string
	User      string
	Password  string

	Insecure bool
)

func init() {
	flag.StringVar(&BaseURL, "url", "", "the base 'Contest URL'")
	flag.StringVar(&ContestId, "c", "", "the 'contest id'")
	flag.StringVar(&User, "u", "", "the 'user'")
	flag.StringVar(&Password, "p", "", "the 'password'")
	flag.BoolVar(&Insecure, "insecure", false, "whether the connection is secure")
}

