package cmd

import (
	ghfactory "github.com/cli/cli/v2/pkg/cmd/factory"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/rafaelcalleja/go-bpkg/pkg/cmd/github"
	"github.com/rafaelcalleja/go-bpkg/pkg/rootcmd"
	"github.com/rafaelcalleja/go-kit/cmd/cobra/version"
	"github.com/rafaelcalleja/go-kit/cmd/helper"
	"github.com/rafaelcalleja/go-kit/cmd/termcolor"
	"github.com/rafaelcalleja/go-kit/logger"
	"github.com/spf13/cobra"
)

type LazyFactory = func() *cmdutil.Factory

func Main() *cobra.Command {
	log := logger.New()
	term := termcolor.NewTermColor()
	errorHelper := helper.NewErrorHelper()
	factory := ghfactory.New(version.GetVersion())

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

	cmd.PersistentFlags().Bool("help", false, "Show help for command")

	cmd.AddCommand(version.NewCmdVersion(errorHelper, log, term))
	cmd.AddCommand(NewPackageInstall(factory, errorHelper, log, term))
	cmd.AddCommand(NewPackageUninstall(errorHelper, log, term))
	cmd.AddCommand(github.NewCmdGithub(factory, errorHelper))

	return cmd
}
