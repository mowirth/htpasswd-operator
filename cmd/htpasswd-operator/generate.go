package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	listener "go.flangaapis.com/htpasswd-kubernetes-operator/htpasswd"
)

var (
	generateCmd = &cobra.Command{
		Use:     "generate",
		Short:   "Run Htpasswd Operator",
		Example: "htpasswd-operator generate",
		Run: func(cmd *cobra.Command, args []string) {
			if err := generate(cmd, args); err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err)
				os.Exit(2)
			}
		},
	}
)

func generate(cmd *cobra.Command, _ []string) error {
	l := listener.HtpasswdListener{
		KubeConfig:   config.KubeConfigFile,
		PasswordFile: config.HTPasswdFile,
		PidFile:      config.PIDFile,
		Namespace:    config.Namespace,
	}

	if err := l.Init(); err != nil {
		return err
	}

	return l.Generate(cmd.Context())
}
