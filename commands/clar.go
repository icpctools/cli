package commands

import (
	"fmt"
	"strings"
	"unicode"

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
	fmt.Printf("\nClarifications (%d):\n", len(clars))

	var table = Table{}
	table.Header = []string{"Time", "Type", "Problem", "Text"}
	table.Align = []int{ALIGN_RIGHT, ALIGN_LEFT, ALIGN_LEFT, ALIGN_LEFT}
	for _, o := range clars {
		var kind = ""
		if o.FromTeamId == "" && o.ToTeamId == "" {
			kind = "Broadcast message from jury"
		} else if o.FromTeamId != "" {
			kind = "Clarification sent to jury"
		} else {
			kind = "Response from jury"
		}
		var prb = ""
		if o.ProblemId != "" {
			problems, err := api.Problems()
			if err != nil {
				return fmt.Errorf("could not get problems; %w", err)
			}

			problem, hasProblem := problemSet(problems).byId(o.ProblemId)
			if !hasProblem {
				prb = "unknown problem"
			} else {
				prb = fmt.Sprintf("%s: %s", problem.Label, problem.Name)
			}
		}
		var first = true
		var time = fmt.Sprintf("%v", o.ContestTime)
		lines := strings.Split(o.Text, "\n")
		for _, s := range lines {
			s = strings.TrimFunc(s, func(r rune) bool {
				return !unicode.IsGraphic(r)
			})
			if first {
				first = false
				table.appendRow([]string{time, kind, prb, s})
			} else {
				table.appendRow([]string{"", "", "", s})
			}
		}
	}
	table.print()

	return nil
}
