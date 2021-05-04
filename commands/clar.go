package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var clarCommand = &cobra.Command{
	Use:     "clar",
	Short:   "Get clarifications",
	RunE:    fetchClars,
	PreRunE: configHelper("baseurl"),
}

func fetchClars(cmd *cobra.Command, args []string) error {
	api, err := contestApi()
	if err != nil {
		return fmt.Errorf("could not connect to the API; %w", err)
	}

	clars, err := api.Clarifications()

	if err != nil {
		return fmt.Errorf("could not retrieve clarifications; %w", err)
	}

	// output
	fmt.Printf("Clarifications (%d):\n", len(clars))
	for _, o := range clars {
		if o.FromTeamId == "" && o.ToTeamId == "" {
			fmt.Printf("  Broadcast message from jury at %v", o.ContestTime)
		} else if o.FromTeamId != "" {
			fmt.Printf("  Clarification sent to jury at %v", o.ContestTime)
		} else {
			fmt.Printf("  Response from jury at %v", o.ContestTime)
		}
		if o.ProblemId != "" {
			problems, err := api.Problems()
			if err != nil {
				return fmt.Errorf("could not get problems; %w", err)
			}

			problem, hasProblem := problemSet(problems).byId(o.ProblemId)
			if !hasProblem {
				fmt.Printf(" for unknown problem")
			}

			fmt.Printf(" for problem %s (%s)", problem.Label, problem.Name)
		}
		fmt.Printf(":\n")
		lines := strings.Split(o.Text, "\n")
		for _, s := range lines {
			fmt.Printf("     %s\n", s)
		}
	}

	return nil
}
