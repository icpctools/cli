package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:                   "login [username] [password]",
		Short:                 "Store login credentials",
		Args:                  cobra.ExactValidArgs(2),
		DisableFlagsInUseLine: true,
		RunE:                  login,
	}

	rootCommand.AddCommand(cmd)
}

func login(cmd *cobra.Command, args []string) error {
	storedConfig.Username = args[0]
	storedConfig.Password = args[1]
	if err := storedConfig.WriteConfig(); err != nil {
		return err
	}

	_, err := fmt.Fprintf(cmd.OutOrStdout(), "Successfully written config to %s\n", storedConfig.ConfigFile())
	return err
}
