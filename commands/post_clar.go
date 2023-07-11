package commands

import (
	"fmt"
	interactor "github.com/icpctools/api-interactor"

	"github.com/spf13/cobra"
)

var postClarCommand = &cobra.Command{
	Use:     "post-clar [text]",
	Short:   "Post a clarification",
	Args:    cobra.ExactValidArgs(1),
	RunE:    postClarification,
	PreRunE: configHelper("baseurl"),
}

func postClarification(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true
	api, err := contestApi()
	if err != nil {
		return fmt.Errorf("could not connect to the server; %w", err)
	}

	if problemId != "" {
		// Get the problems and languages
		problems, err := interactor.List(api, interactor.Problem{})
		if err != nil {
			return fmt.Errorf("could not get problems; %w", err)
		}

		problem, hasProblem := problemSet(problems).byId(problemId)

		if !hasProblem {
			return fmt.Errorf("couldn't find the problem specified")
		}

		problemId = problem.Id
	}

	clar, err := api.PostClarification(problemId, args[0])
	if err != nil {
		return fmt.Errorf("could not post clarification: %w", err)
	}

	_, err = fmt.Fprintf(cmd.OutOrStdout(), "Clarification accepted at %s\n", clar.ContestTime)
	return err
}
