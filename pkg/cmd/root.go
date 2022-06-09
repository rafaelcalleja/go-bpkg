package cmd

import (
	"github.com/rafaelcalleja/go-bpkg/pkg/rootcmd"
	"github.com/rafaelcalleja/go-kit/cmd/cobra/version"
	"github.com/rafaelcalleja/go-kit/cmd/helper"
	"github.com/rafaelcalleja/go-kit/cmd/termcolor"
	"github.com/rafaelcalleja/go-kit/logger"
	"github.com/spf13/cobra"
)

func Main() *cobra.Command {
	log := logger.New()

	cmd := &cobra.Command{
		Use:   rootcmd.TopLevelCommand,
		Short: "Bash Package Manager Go Client",
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			if err != nil {
				log.Errorf(err.Error())
			}
		},
	}

	cmd.AddCommand(version.NewCmdVersion(helper.NewErrorHelper(), log, termcolor.NewTermColor()))

	return cmd
}
