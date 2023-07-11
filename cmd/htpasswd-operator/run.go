package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.flangaapis.com/htpasswd-kubernetes-operator/htpasswd"
)

var (
	runCmd = &cobra.Command{
		Use:     "watch",
		Short:   "Start Htpasswd Operator",
		Example: "htpasswd-operator watch",
		Run: func(cmd *cobra.Command, args []string) {
			if err := run(cmd, args); err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err)
				os.Exit(2)
			}
		},
	}
)

func run(cmd *cobra.Command, _ []string) error {
	l := listener.HtpasswdListener{
		KubeConfig:   config.KubeConfigFile,
		PasswordFile: config.HTPasswdFile,
		PidFile:      config.PIDFile,
		Namespace:    config.Namespace,
	}

	if err := l.Init(); err != nil {
		return err
	}

	l.Run(cmd.Context())
	return nil
}
