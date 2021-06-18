package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var setCommand = &cobra.Command{
	Use:                   "set",
	Short:                 "Set base URL or contest ID",
	DisableFlagsInUseLine: true,
}

var setUrlCommand = &cobra.Command{
	Use:                   "url [url]",
	Short:                 "Set base URL",
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	RunE:                  setUrl,
}

var setIdCommand = &cobra.Command{
	Use:                   "id [id]",
	Short:                 "Set contest ID",
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	RunE:                  setId,
}

func setUrl(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true
	viper.Set("baseurl", args[0])
	if err := viper.WriteConfigAs(configFile()); err != nil {
		return err
	}

	fmt.Printf("Base URL set to %s.\n", args[0])
	return nil
}

func setId(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true
	viper.Set("contest", args[0])
	if err := viper.WriteConfigAs(configFile()); err != nil {
		return err
	}

	fmt.Printf("Contest ID set to %s.\n", args[0])
	return nil
}
