package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	// "github.com/spf13/viper"
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
	problemId string

	insecure bool

	storedConfig = &config.Config{}
)

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
	// Load root command
	rootCommand.PersistentFlags().StringVarP(&baseUrl, "baseurl", "b", "", "base URL to use")
	rootCommand.PersistentFlags().StringVarP(&username, "username", "u", "", "username to communicate with the API")
	rootCommand.PersistentFlags().StringVarP(&password, "password", "p", "", "password to communicate with the API")
	rootCommand.PersistentFlags().StringVarP(&contestId, "contest", "c", "", "contest ID to use")
	rootCommand.PersistentFlags().BoolVarP(&insecure, "insecure", "i", false, "whether to allow insecure HTTPS connections")
	rootCommand.PersistentFlags().StringVar(&problemId, "problem", "", "problem ID to post a clarification for. Leave empty for general clarification")

	rootCommand.Long = fmt.Sprintf(`%s

Note that if the [-b/--baseurl], [-c/--contest], [-i/--insecure], [-p/--password] and [-u/--username] flags
are not supplied, they are read from the configuration file (%s)`, rootCommand.Short, storedConfig.ConfigFile())

	// Bind all valus
	allFlags := []string{"baseurl", "username", "password", "contest", "insecure"}
	for _, flag := range allFlags {
		if err := viper.BindPFlag(flag, rootCommand.PersistentFlags().Lookup(flag)); err != nil {
			// TODO replace this with a better method
			panic(err)
		}
	}

	// Register the contests command
	addSub(rootCommand, contestCommand, "baseurl")
	addSub(rootCommand, postClarCommand, "baseurl", "contest")
	addSub(rootCommand, problemCommand, "baseurl", "contest")
}

func addSub(root, sub *cobra.Command, requiredFlags ...string) {
	sub.PersistentFlags().AddFlagSet(root.PersistentFlags())
	sub.Flags().AddFlagSet(root.Flags())
	root.AddCommand(sub)

	for _, flag := range requiredFlags {
		if err := sub.MarkPersistentFlagRequired(flag); err != nil {
			// TODO better error handling, though this is not recoverable
			panic(err)
		}
	}
}
