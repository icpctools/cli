package commands

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	interactor "github.com/tuupke/api-interactor"
)

func init() {
	cmd := &cobra.Command{
		Use:   "problem",
		Short: "Get problems",
		RunE:  fetchProblems,
	}

	rootCommand.AddCommand(cmd)
}

func fetchProblems(cmd *cobra.Command, args []string) error {
	if baseUrl == "" {
		return errors.New("no base URL provided in flag or config")
	}
	if contestId == "" {
		return errors.New("no contest ID provided in flag or config")
	}

	api, err := interactor.ContestInteractor(baseUrl, username, password, contestId, insecure)
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
