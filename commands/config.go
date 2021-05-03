package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"sort"
	"strings"
)

var configCommand = &cobra.Command{
	Use:                   "config",
	Short:                 "Set base URL or contest ID",
	DisableFlagsInUseLine: true,
	SilenceUsage:          true,
}

var configListCommand = &cobra.Command{
	Use:                   "list",
	Short:                 "List loaded config",
	RunE:                  fetchConfig,
	SilenceUsage:          true,
	DisableFlagsInUseLine: true,
}

var configDeleteCommand = &cobra.Command{
	Use:                   "delete [key]",
	Short:                 "Delete stored config",
	Args:                  cobra.ExactArgs(1),
	RunE:                  deleteConfig,
	SilenceUsage:          true,
	DisableFlagsInUseLine: true,
}

var configSetCommand = &cobra.Command{
	Use:                   "set [key] [value]",
	Short:                 "Set some key in the config",
	Args:                  cobra.ExactArgs(2),
	RunE:                  setConfig,
	SilenceUsage:          true,
	DisableFlagsInUseLine: true,
}

// configKeys returns all keys currently known by Viper that can be set. The keys it returns can optionally be filtered,
// removing more keys from the returned list. The returned keys are sorted before being returned.
func configKeys(ignoreKeys ...string) []string {
	// Construct a map to lookup which keys must be ignored
	var ignoredConfig = map[string]struct{}{"problem": {}}
	for _, i := range ignoreKeys {
		ignoredConfig[i] = struct{}{}
	}

	allKeys := viper.AllKeys()
	all := make([]string, 0, len(allKeys))
	for _, k := range allKeys {
		// Check if this config key must be filtered
		if _, ok := ignoredConfig[k]; ok {
			continue
		}

		all = append(all, k)
	}

	// Sort the remaining keys before returning
	sort.Strings(all)
	return all
}

func fetchConfig(cmd *cobra.Command, args []string) error {
	viper.AllSettings()

	// Determine the amount of padding needed
	var maxLength int
	keys := configKeys()
	for _, k := range keys {
		if maxLength < len(k) {
			maxLength = len(k)
		}
	}

	// Construct a format string that outputs the keys in a padded manner
	format := fmt.Sprintf("  %%-%vv %%v\n", maxLength+6)

	fmt.Println("The current config is:")
	for _, k := range keys {
		fmt.Printf(format, k, viper.Get(k))
	}

	return nil
}

func deleteConfig(cmd *cobra.Command, args []string) error {
	keys := configKeys("username", "password")
	for _, v := range keys {
		if v == args[0] {
			viper.Set(v, nil)
			if err := viper.WriteConfigAs(configFile()); err != nil {
				return fmt.Errorf("could not overwrite config file; %w", err)
			}

			return nil
		}
	}

	// Key not found, throw an error
	formatted, verb := keyFormat(keys)
	return fmt.Errorf("unknown config value to unset: '%s'. Only %s %s allowed", args[0], formatted, verb)
}

func setConfig(cmd *cobra.Command, args []string) error {
	keys := configKeys("username", "password")
	for _, v := range keys {
		if v == args[0] {
			viper.Set(v, args[0])
			if err := viper.WriteConfigAs(configFile()); err != nil {
				return fmt.Errorf("could not overwrite config file; %w", err)
			}

			return nil
		}
	}

	// Key not found, throw an error
	formatted, verb := keyFormat(keys)
	return fmt.Errorf("unknown config value to unset: '%s'. Only %s %s allowed", args[0], formatted, verb)
}

// keyFormat formats a slice of strings (in this case, a set of config flags) in a human readable manner.
func keyFormat(keys []string) (string, string) {
	switch len(keys) {
	case 0:
		return "the empty set", "is"
	case 1:
		return keys[1], "is"
	case 2:
		return fmt.Sprintf("%v and %v", keys[0], keys[1]), "are"
	default:
		last := keys[len(keys)-1]
		return strings.Join(keys[:len(keys)-1], ", ") + ", and " + last, "are"
	}
}
