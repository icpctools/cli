package commands

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
	interactor "github.com/tuupke/api-interactor"
)

var problemCommand = &cobra.Command{
	Use:     "problem",
	Short:   "Get problems",
	RunE:    fetchProblems,
	PreRunE: configHelper("baseurl", "contest"),
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

	_, err = fmt.Fprint(cmd.OutOrStdout(), p)
	return err
}
