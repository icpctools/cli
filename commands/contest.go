package commands

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	interactor "github.com/tuupke/api-interactor"
)

func init() {
	cmd := &cobra.Command{
		Use:   "contest",
		Short: "Get contests",
		RunE:  fetchContests,
	}

	rootCommand.AddCommand(cmd)
}

func fetchContests(cmd *cobra.Command, args []string) error {
	if baseUrl == "" {
		return errors.New("no base URL provided in flag or config")
	}

	api, err := interactor.ContestsInteractor(baseUrl, username, password, insecure)
	if err != nil {
		return fmt.Errorf("could not connect to the API; %w", err)
	}

	var c interface{}
	if contestId != "" {
		c, err = api.ContestById(contestId)
	} else {
		c, err = api.Contests()
	}

	if err != nil {
		return fmt.Errorf("could not retrieve contests; %w", err)
	}

	_, err = fmt.Fprint(cmd.OutOrStdout(), c)
	return err
}
