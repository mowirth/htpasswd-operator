package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"htpasswd-operator/internal/operator"
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
	l := operator.HtpasswdOperator{
		KubeConfig:   config.KubeConfigFile,
		PasswordFile: config.HtpasswdFile,
		PidFile:      config.PIDFile,
	}

	if err := l.Init(cmd.Context()); err != nil {
		return err
	}

	l.Run(cmd.Context())
	return nil
}
