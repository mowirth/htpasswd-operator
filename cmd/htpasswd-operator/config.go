package main

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Configuration struct {
	HTPasswdFile   string `mapstructure:"htpasswd_file"`
	Namespace      string `mapstructure:"namespace"`
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

	if s.HTPasswdFile == "" {
		return fmt.Errorf("htpasswd filepath must not be nil")
	}

	if s.Namespace == "" {
		return fmt.Errorf("namespace must not be empty")
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
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
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
