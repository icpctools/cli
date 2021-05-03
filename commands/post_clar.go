package commands

import (
	"fmt"
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
	api, err := contestApi()
	if err != nil {
		return fmt.Errorf("could not connect to the API; %w", err)
	}

	if problemId != "" {
		// Get the problems and languages
		problems, err := api.Problems()
		if err != nil {
			return fmt.Errorf("could not get problems; %w", err)
		}

		problem, hasProblem := problemSet(problems).byId(problemId)

		if !hasProblem {
			return fmt.Errorf("no known problem specified")
		}

		problemId = problem.Id
	}

	clarId, err := api.PostClarification(problemId, args[0])

	if err != nil {
		return fmt.Errorf("could not post clarification: %w", err)
	}

	_, err = fmt.Fprintf(cmd.OutOrStdout(), "Clarification posted. ID: %s\n", clarId)
	return err
}
