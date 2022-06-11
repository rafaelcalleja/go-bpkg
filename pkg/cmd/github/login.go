package github

import (
	"github.com/cli/cli/v2/pkg/cmd/auth/login"
	"github.com/cli/cli/v2/pkg/cmdutil"
	cmdHelper "github.com/rafaelcalleja/go-kit/cmd/helper"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func NewCmdLogin(cmdFactory *cmdutil.Factory, helper cmdHelper.ErrorHelper) *cobra.Command {
	o := &BaseOptions{}

	githubLoginCmd := login.NewCmdLogin(cmdFactory, nil)

	newCmd := &cobra.Command{
		Use:   "login",
		Short: "Github command line tool",
		Run: func(cmd *cobra.Command, args []string) {
			if "" != strings.TrimSpace(o.configDir) {
				_ = os.Setenv("GH_CONFIG_DIR", o.configDir)
			}

			os.Unsetenv("GITHUB_TOKEN")
			err := githubLoginCmd.RunE(cmd, args)
			helper.CheckErr(err)

			err = login.NewCmdLogin(cmdFactory, func(options *login.LoginOptions) error {
				o.hostname = options.Hostname
				return nil
			}).RunE(cmd, args)
			helper.CheckErr(err)

			cfg, err := cmdFactory.Config()

			if "" == strings.TrimSpace(o.hostname) {
				o.hostname, err = cfg.DefaultHost()
				helper.CheckErr(err)
			}

			_, _, _ = cfg.GetWithSource(o.hostname, "oauth_token")
		},
	}

	newCmd.Flags().AddFlagSet(githubLoginCmd.Flags())

	return newCmd
}
