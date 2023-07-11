package commands

import (
	"fmt"
	interactor "github.com/icpctools/api-interactor"

	"github.com/spf13/cobra"
)

var scoreboardCommand = &cobra.Command{
	Use:     "scoreboard",
	Short:   "Show the contest scoreboard",
	Args:    cobra.NoArgs,
	RunE:    scoreboard,
	PreRunE: configHelper("baseurl"),
}

func scoreboard(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true
	api, err := contestApi()
	if err != nil {
		return fmt.Errorf("could not connect to the server; %w", err)
	}

	problems, err := interactor.List(api, interactor.Problem{})
	if err != nil {
		return fmt.Errorf("could not retrieve problems; %w", err)
	}

	t, err := interactor.List(api, interactor.Team{})
	if err != nil {
		return fmt.Errorf("could not retrieve teams; %w", err)
	}

	sc, err := api.Scoreboard()
	if err != nil {
		return fmt.Errorf("could not retrieve scoreboard; %w", err)
	}

	fmt.Printf("\nContest Scoreboard\n")
	var table = Table{}
	table.Header = []string{"Rank", "Team"}
	table.Align = []int{ALIGN_RIGHT, ALIGN_LEFT}
	for _, p := range problems {
		table.Header = append(table.Header, p.Label)
		table.Align = append(table.Align, ALIGN_RIGHT)
	}
	table.Header = append(table.Header, "Solved", "Time")
	table.Align = append(table.Align, ALIGN_RIGHT, ALIGN_RIGHT)
	for _, r := range sc.Rows {
		team, _ := teamSet(t).byId(string(r.TeamId))
		var name = team.Name
		if team.DisplayName != "" {
			name = team.DisplayName
		}
		var row = []string{fmt.Sprintf("%d", r.Rank), team.Id + ": " + name}

		for _, p := range problems {
			var solved bool
			for _, rp := range r.Problems {
				if rp.Solved && string(rp.ProblemId) == p.Id {
					solved = true
				}
			}
			if solved {
				row = append(row, p.Label)
			} else {
				row = append(row, "")
			}
		}

		row = append(row, fmt.Sprintf("%d", r.Score.NumSolved), fmt.Sprintf("%v", r.Score.TotalTime))
		table.appendRow(row)
	}

	table.print()

	return nil
}
