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
	Short:                 "Store base URL",
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	RunE:                  setUrl,
}

var setIdCommand = &cobra.Command{
	Use:                   "id [id]",
	Short:                 "Store contest Id",
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	RunE:                  setId,
}

func setUrl(cmd *cobra.Command, args []string) error {
	viper.Set("baseurl", args[0])
	if err := viper.WriteConfigAs(configFile()); err != nil {
		return err
	}

	fmt.Printf("Successfully written config to %s\n", configFile())
	return nil
}

func setId(cmd *cobra.Command, args []string) error {
	viper.Set("contest", args[0])
	if err := viper.WriteConfigAs(configFile()); err != nil {
		return err
	}

	fmt.Printf("Successfully written config to %s\n", configFile())
	return nil
}
