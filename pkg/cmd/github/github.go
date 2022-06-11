package github

import (
	"github.com/cli/cli/v2/pkg/cmdutil"
	cmdHelper "github.com/rafaelcalleja/go-kit/cmd/helper"
	"github.com/spf13/cobra"
)

type BaseOptions struct {
	hostname  string
	configDir string
}

func NewCmdGithub(
	cmdFactory *cmdutil.Factory,
	helper cmdHelper.ErrorHelper,
) *cobra.Command {

	cmd := &cobra.Command{
		Use:     "github <command>",
		Aliases: []string{"gh"},
		Short:   "Login, logout, and refresh your authentication",
		Long:    `Manage gh's authentication state.`,
	}

	cmdutil.DisableAuthCheck(cmd)

	cmd.AddCommand(NewCmdLogin(cmdFactory, helper))
	cmd.AddCommand(NewCmdStatus(cmdFactory, helper))

	return cmd
}
