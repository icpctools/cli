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
	Short: "Get contests",
	RunE:  fetchContests,
}

func outputContest(c interactor.Contest) {
	fmt.Printf(" %10s: %s\n", c.Id, c.Name)
	if c.StartTime == (interactor.ApiTime{}) {
		if c.CountdownTime != interactor.ApiRelTime(0) {
			fmt.Printf("             %v (countdown paused at %s)\n", c.Duration, c.CountdownTime)
		} else {
			fmt.Printf("             %v (not scheduled)\n", c.Duration)
		}
	} else {
		now := time.Now()
		if c.StartTime.Time().After(now) {
			fmt.Printf("             %v starting at %v\n", c.Duration, c.StartTime)
		} else if (c.StartTime.Time().Add(c.Duration.Duration())).After(now) {
			fmt.Printf("             Contest running. %v started at %v\n", c.Duration, c.StartTime)
		} else {
			fmt.Printf("             Contest over. Started at %v\n", c.StartTime)
		}
	}
}

func fetchContests(cmd *cobra.Command, args []string) error {
	api, err := contestsApi()
	if err != nil {
		return fmt.Errorf("could not connect to the API; %w", err)
	}

	if contestId != "" {
		c, err := api.ContestById(contestId)

		if err != nil {
			return fmt.Errorf("could not retrieve contest; %w", err)
		}

		outputContest(c)
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
		fmt.Printf("Contests (%d):\n", len(c))
		for _, o := range c {
			outputContest(o)
		}
	}

	return nil
}
