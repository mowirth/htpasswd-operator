package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Configuration struct {
	HtpasswdFile   string `mapstructure:"htpasswd_file"`
	KubeConfigFile string `mapstructure:"kube_config"`
	PIDFile        string `mapstructure:"pid_file"`
	LogLevel       string `mapstructure:"log_level"`
}

func (s *Configuration) Validate() error {
	level, err := logrus.ParseLevel(s.LogLevel)
	if err != nil {
		return err
	}

	logrus.SetLevel(level)

	if s.HtpasswdFile == "" {
		return fmt.Errorf("operator filepath must not be nil")
	}

	return nil
}

func (s *Configuration) Read() error {
	if s == nil {
		return fmt.Errorf("empty config")
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.SetDefault("log_level", "INFO")
	viper.SetDefault("htpasswd_file", "./htpasswd")

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		switch {
		case errors.As(err, &configFileNotFoundError):
		default:
			return err
		}
	}

	bindEnvs(*s)
	if err := viper.Unmarshal(s); err != nil {
		return err
	}

	return s.Validate()
}
