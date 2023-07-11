package commands

import (
	"fmt"
	interactor "github.com/icpctools/api-interactor"
	"sort"

	"github.com/spf13/cobra"
)

var problemCommand = &cobra.Command{
	Use:     "problem",
	Short:   "List problems",
	Args:    cobra.NoArgs,
	RunE:    fetchProblems,
	PreRunE: configHelper("baseurl"),
}

func fetchProblems(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true
	api, err := contestApi()
	if err != nil {
		return fmt.Errorf("could not connect to the server; %w", err)
	}

	p, err := interactor.List(api, interactor.Problem{})
	if err != nil {
		return fmt.Errorf("could not retrieve problems; %w", err)
	}

	// sort by ordinal
	sort.Slice(p, func(i, j int) bool {
		return p[i].Ordinal < p[j].Ordinal
	})

	// output
	fmt.Printf("\nProblems (%d):\n", len(p))

	var table = Table{}
	table.Header = []string{"Label", "Name"}
	table.Align = []int{ALIGN_LEFT, ALIGN_LEFT}
	for _, o := range p {
		table.appendRow([]string{o.Label, o.Name})
	}
	table.print()

	return nil
}
