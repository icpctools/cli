package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var loginCommand = &cobra.Command{
	Use:                   "login [username] [password]",
	Short:                 "Set login credentials",
	Args:                  cobra.ExactValidArgs(2),
	DisableFlagsInUseLine: true,
	RunE:                  login,
}

func login(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true
	viper.Set("username", args[0])
	viper.Set("password", args[1])
	if err := viper.WriteConfigAs(configFile()); err != nil {
		return err
	}

	fmt.Printf("Credentials set for %s\n", args[0])
	return nil
}
