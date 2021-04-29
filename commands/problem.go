package commands

import (
	"fmt"
	"sort"
	"github.com/spf13/cobra"
)

var problemCommand = &cobra.Command{
	Use:     "problem",
	Short:   "Get problems",
	RunE:    fetchProblems,
	PreRunE: configHelper("baseurl", "contest"),
}

func fetchProblems(cmd *cobra.Command, args []string) error {
	api, err := contestApi()
	if err != nil {
		return fmt.Errorf("could not connect to the API; %w", err)
	}

	p, err := api.Problems()

	if err != nil {
		return fmt.Errorf("could not retrieve problems; %w", err)
	}

	// sort by ordinal
 	sort.Slice(p, func(i, j int) bool {
 		return p[i].Ordinal < p[j].Ordinal
 	})

 	// output
 	fmt.Printf("Problems (%d):\n", len(p))
 	for _, o := range p {
 		fmt.Printf(" %3s: %s\n", o.Label, o.Name)
 	}

 	return nil
}
