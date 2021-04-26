package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"tools.icpc.global/contest/config"
)

var (
	rootCommand = &cobra.Command{
		Use:              "contest",
		Short:            "A CLI tool for CCS Api access",
		PersistentPreRun: mergeConfig,
	}

	baseUrl   string
	username  string
	password  string
	contestId string

	insecure bool

	storedConfig = &config.Config{}
)

func init() {
	rootCommand.Long = fmt.Sprintf(`%s

Note that if the [-b/--baseurl], [-c/--contest], [-i/--insecure], [-p/--password] and [-u/--username] flags
are not supplied, they are read from the configuration file (%s)`, rootCommand.Short, storedConfig.ConfigFile())
}

func Execute() error {
	return rootCommand.Execute()
}

func mergeConfig(_ *cobra.Command, _ []string) {
	// Read config from disk
	storedConfig.ReadConfig()

	// Now merge it into our local config, if local config is empty
	if baseUrl == "" {
		baseUrl = storedConfig.BaseUrl
	}
	if username == "" {
		username = storedConfig.Username
	}
	if password == "" {
		password = storedConfig.Password
	}
	if contestId == "" {
		contestId = storedConfig.ContestId
	}
	if !insecure {
		insecure = storedConfig.Insecure
	}
}

func init() {
	rootCommand.PersistentFlags().StringVarP(&baseUrl, "baseurl", "b", "", "base URL to use")
	rootCommand.PersistentFlags().StringVarP(&username, "username", "u", "", "username to communicate with the API")
	rootCommand.PersistentFlags().StringVarP(&password, "password", "p", "", "password to communicate with the API")
	rootCommand.PersistentFlags().StringVarP(&contestId, "contest", "c", "", "contest ID to use")
	rootCommand.PersistentFlags().BoolVarP(&insecure, "insecure", "i", false, "whether to allow insecure HTTPS connections")
}
