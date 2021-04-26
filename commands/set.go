package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	baseCmd := &cobra.Command{
		Use:                   "set",
		Short:                 "Set base URL or contest ID",
		DisableFlagsInUseLine: true,
	}

	rootCommand.AddCommand(baseCmd)

	baseCmd.AddCommand(&cobra.Command{
		Use:                   "url [url]",
		Short:                 "Store base URL",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		RunE:                  setUrl,
	})
	baseCmd.AddCommand(&cobra.Command{
		Use:                   "id [id]",
		Short:                 "Store contest Id",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		RunE:                  setId,
	})
}

func setUrl(cmd *cobra.Command, args []string) error {
	storedConfig.BaseUrl = args[0]
	if err := storedConfig.WriteConfig(); err != nil {
		return err
	}

	_, err := fmt.Fprintf(cmd.OutOrStdout(), "Successfully written config to %s\n", storedConfig.ConfigFile())
	return err
}

func setId(cmd *cobra.Command, args []string) error {
	storedConfig.ContestId = args[0]
	if err := storedConfig.WriteConfig(); err != nil {
		return err
	}

	_, err := fmt.Fprintf(cmd.OutOrStdout(), "Successfully written config to %s\n", storedConfig.ConfigFile())
	return err
}
