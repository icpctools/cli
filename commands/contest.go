package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

var contestCommand = &cobra.Command{
	Use:   "contest",
	Short: "Get contests",
	RunE:  fetchContests,
}

func fetchContests(cmd *cobra.Command, args []string) error {
	api, err := contestsApi()
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

	fmt.Println(c)
	return nil
}
