package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

type rowStr []string

const ALIGN_LEFT = 0
const ALIGN_RIGHT = 1

type Table struct {
	Header rowStr
	Rows   []rowStr
	Align  []int
}

var scoreboardCommand = &cobra.Command{
	Use:     "scoreboard",
	Short:   "Show the contest scoreboard",
	RunE:    scoreboard,
	PreRunE: configHelper("baseurl"),
}

func (table Table) print() error {
	// determine the amount of padding needed
	var numCol = len(table.Header)
	var maxLength []int = make([]int, numCol)
	var format []string = make([]string, numCol)

	// find max header width
	for i, s := range table.Header {
		if maxLength[i] < len(s) {
			maxLength[i] = len(s)
		}
	}

	// find max cell width
	for _, r := range table.Rows {
		for i, s := range r {
			if maxLength[i] < len(s) {
				maxLength[i] = len(s)
			}
		}
	}

	// create format for each column, respecting width and alignment
	for i := range table.Header {
		if table.Align[i] == ALIGN_LEFT {
			format[i] = fmt.Sprintf("  %%-%vv", maxLength[i])
		} else {
			format[i] = fmt.Sprintf("  %%%vv", maxLength[i])
		}
	}

	// output header bold and underlined
	fmt.Printf("\033[1;4m")
	for i, k := range table.Header {
		fmt.Printf(format[i], k)
	}
	fmt.Printf("\033[0m\n")

	// output each cell
	for _, r := range table.Rows {
		for i, s := range r {
			fmt.Printf(format[i], s)
		}
		fmt.Printf("\n")
	}

	return nil
}

func scoreboard(cmd *cobra.Command, args []string) error {
	api, err := contestApi()
	if err != nil {
		return fmt.Errorf("could not connect to the API; %w", err)
	}

	problems, err := api.Problems()
	if err != nil {
		return fmt.Errorf("could not retrieve problems; %w", err)
	}

	t, err := api.Teams()
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
		var row = []string{fmt.Sprintf("%d", r.Rank), team.Id + ": " + team.Name}

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
		table.Rows = append(table.Rows, row)
	}

	table.print()

	return nil
}
