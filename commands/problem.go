package commands

import (
	"errors"
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	interactor "github.com/tuupke/api-interactor"
)

var problemCommand = &cobra.Command{
	Use:   "problem",
	Short: "Get problems",
	RunE:  fetchProblems,
}

func fetchProblems(cmd *cobra.Command, args []string) error {
	if viper.GetString("baseurl") == "" {
		return errors.New("no base URL provided in flag or config")
	}
	if viper.GetString("contest") == "" {
		return errors.New("no contest ID provided in flag or config")
	}

	api, err := interactor.ContestInteractor(viper.GetString("baseurl"), viper.GetString("username"), viper.GetString("password"), viper.GetString("contest"), viper.GetBool("insecure"))
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
