package commands

import (
	"github.com/spf13/cobra"
)

var (
	rootCommand = &cobra.Command{
		Use:   "contest",
		Short: "A CLI tool for CCS Api access",
	}

	baseUrl   string
	username  string
	password  string
	contestId string

	insecure bool
)

func Execute() error {
	return rootCommand.Execute()
}

func init() {
	rootCommand.PersistentFlags().StringVarP(&baseUrl, "baseurl", "b", "", "")
	if err := rootCommand.MarkPersistentFlagRequired("baseurl"); err != nil {
		panic(err)
	}

	rootCommand.PersistentFlags().StringVarP(&username, "user", "u", "", "")
	rootCommand.PersistentFlags().StringVarP(&password, "pass", "p", "", "")
	rootCommand.PersistentFlags().StringVarP(&contestId, "contest", "c", "", "")
	rootCommand.PersistentFlags().BoolVarP(&insecure, "insecure", "i", false, "")
}
