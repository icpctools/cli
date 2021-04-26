package commands

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	interactor "github.com/tuupke/api-interactor"
)

func init() {
	cmd := &cobra.Command{
		Use:   "post-clar [text]",
		Short: "Post a clarification",
		Args:  cobra.ExactValidArgs(1),
		RunE:  postClarification,
	}

	cmd.Flags().StringVar(&problemId, "problem", "", "problem ID to post a clarification for. Leave empty for general clarification")

	rootCommand.AddCommand(cmd)
}

func postClarification(cmd *cobra.Command, args []string) error {
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

	clarId, err := api.PostClarification(problemId, args[0])

	if err != nil {
		return fmt.Errorf("could not post clarification: %w", err)
	}

	_, err = fmt.Fprintf(cmd.OutOrStdout(), "Clarification posted. ID: %s\n", clarId)
	return err
}
