package commands

import (
	"fmt"
	"github.com/kirsle/configdir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"path/filepath"
)

const (
	configFolder = "icpc"
	configName   = "settings"
	configType   = "yaml"
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
	problemId string

	insecure bool
)

func Execute() error {
	return rootCommand.Execute()
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
are not supplied, they are read from the configuration file (%s)`, rootCommand.Short, configFile())

	// Set viper path and file
	viper.AddConfigPath(configdir.LocalConfig(configFolder))
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)

	// Bind all values
	allFlags := []string{"baseurl", "username", "password", "contest", "insecure"}
	for _, flag := range allFlags {
		if err := viper.BindPFlag(flag, rootCommand.PersistentFlags().Lookup(flag)); err != nil {
			// TODO replace this with a better method
			panic(err)
		}
	}

	// Read in viper config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// TODO replace this with a better method
			panic(err)
		}
	}

	// Register the subcommands
	addSub(rootCommand, contestCommand, "baseurl")
	addSub(rootCommand, postClarCommand, "baseurl", "contest")
	addSub(rootCommand, problemCommand, "baseurl", "contest")
	addSub(rootCommand, loginCommand)
	addSub(rootCommand, logoutCommand)
	addSub(rootCommand, setCommand)
	addSub(setCommand, setUrlCommand)
	addSub(setCommand, setIdCommand)
}

func addSub(root, sub *cobra.Command, requiredFlags ...string) {
	// TODO: contest and baseurl are now required for all commands, not only the three we want them to have

	// TODO: contest and baseurl are required even if they appear in the config file. This should not be the case

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

func configFile() string {
	return fmt.Sprintf("%s.%s", filepath.Join(configdir.LocalConfig(configFolder), configName), configType)
}
