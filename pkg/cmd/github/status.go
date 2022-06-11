package github

import (
	"github.com/cli/cli/v2/pkg/cmd/auth/status"
	"github.com/cli/cli/v2/pkg/cmdutil"
	cmdHelper "github.com/rafaelcalleja/go-kit/cmd/helper"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func NewCmdStatus(cmdFactory *cmdutil.Factory, helper cmdHelper.ErrorHelper) *cobra.Command {
	o := &BaseOptions{}
	githubStatusCmd := status.NewCmdStatus(cmdFactory, nil)

	newCmd := &cobra.Command{
		Use:   "status",
		Short: "Github Auth Status command line tool",
		Run: func(cmd *cobra.Command, args []string) {
			if "" != strings.TrimSpace(o.configDir) {
				_ = os.Setenv("GH_CONFIG_DIR", o.configDir)
			}

			err := githubStatusCmd.RunE(cmd, args)
			helper.CheckErr(err)
		},
	}

	newCmd.Flags().AddFlagSet(githubStatusCmd.Flags())
	newCmd.Flags().StringVar(&o.configDir, "config", "", "Config path")

	return newCmd
}
