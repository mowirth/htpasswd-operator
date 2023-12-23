package main

import (
	"fmt"
	"os"

	"github.com/mowirth/htpasswd-operator/internal/operator"
	"github.com/spf13/cobra"
)

var (
	runCmd = &cobra.Command{
		Use:     "watch",
		Short:   "Start Htpasswd Operator",
		Example: "github.com/mowirth/htpasswd-operator watch",
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
