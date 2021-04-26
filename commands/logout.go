package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:                   "logout",
		Short:                 "Remove login credentials",
		DisableFlagsInUseLine: true,
		RunE:                  logout,
	}

	rootCommand.AddCommand(cmd)
}

func logout(cmd *cobra.Command, args []string) error {
	storedConfig.Username = ""
	storedConfig.Password = ""
	if err := storedConfig.WriteConfig(); err != nil {
		return err
	}

	_, err := fmt.Fprintf(cmd.OutOrStdout(), "Successfully written config to %s\n", storedConfig.ConfigFile())
	return err
}
