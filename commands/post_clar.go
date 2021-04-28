package commands

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var postClarCommand = &cobra.Command{
	Use:     "post-clar [text]",
	Short:   "Post a clarification",
	Args:    cobra.ExactValidArgs(1),
	RunE:    postClarification,
	PreRunE: configHelper("baseurl", "contest"),
}

func postClarification(cmd *cobra.Command, args []string) error {
	if viper.GetString("baseurl") == "" {
		return errors.New("no base URL provided in flag or config")
	}
	if viper.GetString("contest") == "" {
		return errors.New("no contest ID provided in flag or config")
	}

	api, err := contestApi()
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
