// Package config provides functions for reading the per-user timetrace config.
package config

import (
	"github.com/spf13/viper"
)

type UseDecimalHours int

const (
	Off UseDecimalHours = iota
	On
	Both
)

type Config struct {
	Store           string          `json:"store"`
	Use12Hours      bool            `json:"use12hours"`
	UseDecimalHours UseDecimalHours `json:"usedecimalhours"`
	Editor          string          `json:"editor"`
	ReportPath      string          `json:"report-path"`
}

var cached *Config

// FromFile reads a configuration file called config.yml and returns it as a
// Config instance. If no configuration file is found, nil and no error will be
// returned. The configuration must live in one of the following directories:
//
//	- /etc/timetrace
//	- $HOME/.timetrace
//	- .
//
// In case multiple configuration files are found, the one in the most specific
// or "closest" directory will be preferred.
func FromFile() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/timetrace/")
	viper.AddConfigPath("$HOME/.timetrace")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	cached = &config

	return cached, nil
}

// Get returns the parsed configuration. The fields of this configuration either
// contain values specified by the user or the zero value of the respective data
// type, e.g. "" for an un-configured string.
//
// Using Get over FromFile avoids the config file from being parsed each time
// the config is needed.
func Get() *Config {
	if cached != nil {
		return cached
	}

	config, err := FromFile()
	if err != nil {
		return &Config{}
	}

	return config
}
