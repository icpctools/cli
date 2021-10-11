package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"time"

	interactor "github.com/icpctools/api-interactor"
	"github.com/kirsle/configdir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configFolder = "icpc"
	configName   = "settings"
	configType   = "yaml"
)

var (
	rootCommand = &cobra.Command{
		Use:   "contest",
		Short: "A CLI tool for CCS Api access",
	}

	baseUrl    string
	username   string
	password   string
	contestId  string
	problemId  string
	languageId string
	entryPoint string

	force    bool
	insecure bool
)

func Execute() error {
	return rootCommand.Execute()
}

func init() {
	// Load root command
	rootCommand.PersistentFlags().StringVarP(&baseUrl, "baseurl", "b", "", "base URL to use")
	rootCommand.PersistentFlags().StringVarP(&username, "username", "u", "", "username to communicate with the API")
	rootCommand.PersistentFlags().StringVarP(&password, "password", "p", "", "password to communicate with the API")
	rootCommand.PersistentFlags().StringVarP(&contestId, "contest", "c", "", "contest ID to use")
	rootCommand.PersistentFlags().BoolVarP(&insecure, "insecure", "i", false, "whether to allow insecure HTTPS connections")

	// Command specific flags
	postClarCommand.Flags().StringVar(&problemId, "problem", "", "problem ID to post a clarification for. Leave empty for general clarification")

	submitCommand.Flags().StringVar(&problemId, "problem", "", "problem ID to submit for. Leave empty to auto detect from first file")
	submitCommand.Flags().StringVarP(&languageId, "language", "l", "", "language ID to submit for. Leave empty to auto detect from first file")
	submitCommand.Flags().StringVarP(&entryPoint, "entry-point", "e", "", "entry point to use. Leave empty if not needed or to auto detect")
	submitCommand.Flags().BoolVarP(&force, "force", "f", false, "whether to force submission (i.e. not ask for confirmation")

	rootCommand.Long = fmt.Sprintf(`%s

Note that if the [-b/--baseurl], [-c/--contest], [-i/--insecure], [-p/--password] and [-u/--username] flags
are not supplied, they are read from the configuration file (%s)`, rootCommand.Short, configFile())

	// Set viper path and file
	configDir := configdir.LocalConfig(configFolder)

	// Ensure config path exists
	err := configdir.MakePath(configDir)
	if err != nil {
		fmt.Printf("can not create config folder: %s\n", err)
		os.Exit(1)
	}

	viper.AddConfigPath(configDir)
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)

	// Bind all values
	allFlags := []string{"baseurl", "username", "password", "contest", "insecure"}
	for _, flag := range allFlags {
		if err := viper.BindPFlag(flag, rootCommand.PersistentFlags().Lookup(flag)); err != nil {
			// TODO replace this with a better method
			panic(err)
		}
	}

	// Read in viper config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// TODO replace this with a better method
			panic(err)
		}
	}

	// Register the subcommands
	setCommand.AddCommand(setUrlCommand)
	setCommand.AddCommand(setIdCommand)
	rootCommand.AddCommand(contestCommand)
	rootCommand.AddCommand(clarCommand)
	rootCommand.AddCommand(postClarCommand)
	rootCommand.AddCommand(problemCommand)
	rootCommand.AddCommand(loginCommand)
	rootCommand.AddCommand(logoutCommand)
	rootCommand.AddCommand(setCommand)
	rootCommand.AddCommand(submitCommand)
	rootCommand.AddCommand(submissionsCommand)
	rootCommand.AddCommand(scoreboardCommand)
}

// configHelper can be used to register which flags must exist. An error is thrown when a required flag is not present
// or set in through viper. If a flag is provided it will override the value stored in viper, such that its interface
// can be used to retrieve all config.
func configHelper(requiredFlags ...string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		for _, flag := range requiredFlags {
			f := cmd.Flag(flag)
			if f != nil && f.Changed {
				viper.Set(flag, f.Value)
				continue
			}

			if viper.IsSet(flag) && !reflect.ValueOf(viper.Get(flag)).IsZero() {
				continue
			}

			// Neither flag, nor value exists, exiting
			return fmt.Errorf("missing flag: '--%v'", flag)
		}

		return nil
	}
}

func configFile() string {
	return fmt.Sprintf("%s.%s", filepath.Join(configdir.LocalConfig(configFolder), configName), configType)
}

// contestApi attempts to load a interactor.ContestApi from the config currently stored in viper.
func contestApi() (interactor.ContestApi, error) {
	contest := viper.GetString("contest")
	if contest == "" {
		api, err := contestsApi()
		if err != nil {
			return nil, fmt.Errorf("could not connect to the API; %w", err)
		}

		c, err := api.Contests()
		if err != nil {
			return nil, fmt.Errorf("could not retrieve contests; %w", err)
		}

		best, err := contestSet(c).bestContest()
		if err != nil {
			return nil, fmt.Errorf("could not pick the best contest; %w", err)
		} else {
			fmt.Printf("Automatically connecting to contest: %s\n", best.Name)
			contest = best.Id
		}
	}
	return interactor.ContestInteractor(
		viper.GetString("baseurl"),
		viper.GetString("username"),
		viper.GetString("password"),
		contest,
		viper.GetBool("insecure"),
	)
}

// contestsApi attempts to load a interactor.ContestsApi from the config currently stored in viper.
func contestsApi() (interactor.ContestsApi, error) {
	return interactor.ContestsInteractor(
		viper.GetString("baseurl"),
		viper.GetString("username"),
		viper.GetString("password"),
		viper.GetBool("insecure"),
	)
}

type (
	contestSet []interactor.Contest
)

func (c contestSet) bestContest() (interactor.Contest, error) {
	var best interactor.Contest

	// simple cases - no contests or only one
	if len(c) == 0 {
		return best, fmt.Errorf("no contests found")
	}
	if len(c) == 1 {
		return c[0], nil
	}

	// ok, there are at least two contests
	// if there is only one contest running, pick it
	var count int
	var unscheduled int
	now := time.Now()
	for _, contest := range c {
		if contest.StartTime == (interactor.ApiTime{}) {
			unscheduled++
		}
		if contest.StartTime.Time().Before(now) && contest.StartTime.Time().Add(time.Duration(contest.Duration)).After(now) {
			best = contest
			count++
		}
	}

	if count == 1 {
		return best, nil
	} else if count >= 2 {
		return interactor.Contest{}, errors.New("more than one contest is currently running")
	}
	if unscheduled == len(c) {
		return interactor.Contest{}, errors.New("there are no scheduled contests")
	}

	// if there is only one contest today, pick it
	count = 0
	for _, contest := range c {
		if dateEqual(contest.StartTime.Time(), now) {
			best = contest
			count++
		}
	}

	if count == 1 {
		return best, nil
	}

	// if all contests are in the future, pick the first one
	sort.Slice(c, func(i, j int) bool {
		return c[i].StartTime.Time().Before(c[j].StartTime.Time())
	})
	if c[0].StartTime != (interactor.ApiTime{}) && c[0].StartTime.Time().After(now) {
		return c[0], nil
	}

	// ok, so all contests are done. just pick the last one
	return c[len(c)-1], nil
}

func dateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
