package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var logoutCommand = &cobra.Command{
	Use:                   "logout",
	Short:                 "Remove login credentials",
	DisableFlagsInUseLine: true,
	RunE:                  logout,
}

func logout(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true
	viper.Set("username", "")
	viper.Set("password", "")
	if err := viper.WriteConfigAs(configFile()); err != nil {
		return err
	}

	fmt.Printf("Successfully written config to %s\n", configFile())
	return nil
}
