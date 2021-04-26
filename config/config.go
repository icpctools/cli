package config

import (
	"fmt"
	"github.com/kirsle/configdir"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

const (
	configFolder = "icpc"
	configFile   = "settings.yaml"
)

type (
	Config struct {
		BaseUrl   string `yaml:"url"`
		ContestId string `yaml:"contest"`
		Username  string `yaml:"username"`
		Password  string `yaml:"password"`
		Insecure  bool   `yaml:"insecure"`
	}
)

func (c *Config) ConfigFile() string {
	configPath := configdir.LocalConfig(configFolder)
	return filepath.Join(configPath, configFile)
}

func (c *Config) ReadConfig() {
	configFile := c.ConfigFile()

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Config file does not exist, do not read anything
		return
	}

	fh, err := os.Open(configFile)
	if err != nil {
		fmt.Printf("Can not read config file: %s\n", err)
		return
	}

	defer fh.Close()

	err = yaml.NewDecoder(fh).Decode(c)
	if err != nil {
		fmt.Printf("Can not read config file: %s\n", err)
	}
}

func (c *Config) WriteConfig() error {
	configPath := configdir.LocalConfig(configFolder)
	// Ensure config path exists
	err := configdir.MakePath(configPath)
	if err != nil {
		return fmt.Errorf("can not create config folder: %s\n", err)
	}

	configFile := c.ConfigFile()

	fh, err := os.Create(configFile)
	if err != nil {
		return fmt.Errorf("can not write config file: %s\n", err)
	}

	defer fh.Close()

	err = yaml.NewEncoder(fh).Encode(c)
	if err != nil {
		return fmt.Errorf("can not write config file: %s\n", err)
	}

	return nil
}
