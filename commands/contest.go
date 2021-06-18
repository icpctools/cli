package commands

import (
	"fmt"
	"sort"
	"time"

	"github.com/spf13/cobra"
	interactor "github.com/tuupke/api-interactor"
)

var contestCommand = &cobra.Command{
	Use:   "contest",
	Short: "List contests",
	RunE:  fetchContests,
}

func outputContest(table *Table, c interactor.Contest) {
	var name = c.FormalName
	if name == "" {
		name = c.Name
	}
	var row = []string{c.Id, name}

	var duration = fmt.Sprintf("%v", c.Duration)
	if c.StartTime == (interactor.ApiTime{}) {
		if c.CountdownTime != interactor.ApiRelTime(0) {
			row = append(row, "", duration, fmt.Sprintf("Countdown paused at %v", c.CountdownTime))
		} else {
			row = append(row, "", duration, "Not scheduled")
		}
	} else {
		now := time.Now()
		if c.StartTime.Time().After(now) {
			row = append(row, fmt.Sprintf("%v", c.StartTime), duration, "Scheduled")
		} else if (c.StartTime.Time().Add(c.Duration.Duration())).After(now) {
			row = append(row, fmt.Sprintf("%v", c.StartTime), duration, "Started")
		} else {
			row = append(row, fmt.Sprintf("%v", c.StartTime), duration, "Contest over")
		}
	}
	table.appendRow(row)
}

func fetchContests(cmd *cobra.Command, args []string) error {
	api, err := contestsApi()
	if err != nil {
		return fmt.Errorf("could not connect to the API; %w", err)
	}

	var table = Table{}
	table.Header = []string{"Id", "Name", "Start Time", "Length", "Status"}
	table.Align = []int{ALIGN_LEFT, ALIGN_LEFT, ALIGN_LEFT, ALIGN_RIGHT, ALIGN_LEFT}
	if contestId != "" {
		c, err := api.ContestById(contestId)

		if err != nil {
			return fmt.Errorf("could not retrieve contest; %w", err)
		}

		outputContest(&table, c)
		table.print()
	} else {
		c, err := api.Contests()

		if err != nil {
			return fmt.Errorf("could not retrieve contests; %w", err)
		}

		// sort by start time
		sort.Slice(c, func(i, j int) bool {
			return c[i].StartTime.Time().Before(c[j].StartTime.Time())
		})

		// output
		fmt.Printf("\nContests (%d):\n", len(c))
		for _, o := range c {
			outputContest(&table, o)
		}
		table.print()
	}

	return nil
}
