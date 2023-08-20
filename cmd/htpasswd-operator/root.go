package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var config Configuration
var (
	rootCmd = &cobra.Command{
		Use: "htpasswd-operator",
		Run: func(cmd *cobra.Command, args []string) {
			if err := run(cmd, args); err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(2)
			}
		},
	}
)

func init() {
	cobra.OnInitialize(func() {
		config = Configuration{}
		if err := config.Read(); err != nil {
			log.Panicf("Failed to initialize config: %v", err)
		}
	})

	rootCmd.AddCommand(runCmd)
}

// Copied from some github issue which I can't find anymore unfortunately (somewhere around spf13/viper).
// Fixes some issues with the env binding of Viper...
func bindEnvs(iface interface{}, parts ...string) {
	ifv := reflect.ValueOf(iface)
	ift := reflect.TypeOf(iface)
	for i := 0; i < ift.NumField(); i++ {
		v := ifv.Field(i)
		t := ift.Field(i)
		tv, ok := t.Tag.Lookup("mapstructure")
		if !ok {
			tv = strings.ToLower(t.Name)
		}
		if tv == "-" {
			continue
		}

		switch v.Kind() {
		case reflect.Struct:
			bindEnvs(v.Interface(), append(parts, tv)...)
		default:
			// Bash doesn't allow env variable names with a dot so
			// bind the double underscore version.
			keyDot := strings.Join(append(parts, tv), ".")
			keyUnderscore := strings.Join(append(parts, tv), "_")
			if err := viper.BindEnv(keyDot, strings.ToUpper(keyUnderscore)); err != nil {
				logrus.Errorf("Failed to bind %v to %v: %v", keyDot, strings.ToUpper(keyUnderscore), err)
			}
		}
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}
}
