package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

var problemCommand = &cobra.Command{
	Use:     "problem",
	Short:   "Get problems",
	RunE:    fetchProblems,
	PreRunE: configHelper("baseurl", "contest"),
}

func fetchProblems(cmd *cobra.Command, args []string) error {
	api, err := contestApi()
	if err != nil {
		return fmt.Errorf("could not connect to the API; %w", err)
	}

	p, err := api.Problems()

	if err != nil {
		return fmt.Errorf("could not retrieve problems; %w", err)
	}

	_, err = fmt.Fprint(cmd.OutOrStdout(), p)
	return err
}
