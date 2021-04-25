package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	interactor "github.com/tuupke/api-interactor"
)

var problemsCommand = &cobra.Command{
	Use:   "problem",
	Short: "Get problemss",
	RunE:  fetchProblems,
}

func init() {

	if err := problemsCommand.MarkPersistentFlagRequired("contest"); err != nil {
		panic(err)
	}

	rootCommand.AddCommand(problemsCommand)


}

func fetchProblems(cmd *cobra.Command, args []string) error {
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
