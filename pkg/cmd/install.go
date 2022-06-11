package cmd

import (
	"fmt"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/rafaelcalleja/go-bpkg/pkg/repository"
	"github.com/rafaelcalleja/go-kit/cmd/helper"
	"github.com/rafaelcalleja/go-kit/cmd/termcolor"
	"github.com/rafaelcalleja/go-kit/logger"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

type PackageInstallOptions struct {
	packageName string
	installPath string
	token       string
	upgrade     bool
	private     bool
}

func NewPackageInstall(
	factory *cmdutil.Factory,
	helper helper.ErrorHelper,
	log logger.Logger,
	term termcolor.TermColor,
) *cobra.Command {
	o := &PackageInstallOptions{}

	newCmd := &cobra.Command{
		Use:   "install",
		Short: "BPKG install",
		Run: func(cmd *cobra.Command, args []string) {
			if "" != strings.TrimSpace(o.token) {
				_ = os.Setenv("GITHUB_TOKEN", o.token)
			}

			fqpVO, err := repository.NewFullyQualifyPackage(o.packageName)
			helper.CheckErr(err)

			releaseVersion, err := repository.NewReleaseVersionWith(
				repository.ReleaseVersionWithOrganization(fqpVO.Organization()),
				repository.ReleaseVersionWithName(fqpVO.Name()),
				repository.ReleaseVersionWithVersion(fqpVO.Version()),
			)

			if fqpVO.Version() == "latest" {
				releaseVersion, err = repository.NewReleaseLatestVersion(
					fqpVO.Organization(),
					fqpVO.Name(),
					repository.NewGithubVersionFinder(factory),
				)
			}

			if true == releaseVersion.IsVersionInstalled(releaseVersion.Version(), o.installPath) {
				log.Infof("Package %s already at version %s", term.ColorInfo(fmt.Sprintf("%s/%s", releaseVersion.Organization, releaseVersion.Name)),
					term.ColorInfo(releaseVersion.Version()))

				return
			}

			log.Infof("Installing Package %s at %s", term.ColorInfo(releaseVersion.String()),
				term.ColorInfo(o.installPath))

			assetGithub := repository.NewGithubProvider(factory)
			err = releaseVersion.DownloadAsset(assetGithub, o.installPath)
			helper.CheckErr(err)

			log.Infof("Installed Successfully")
		},
	}

	newCmd.Flags().StringVar(&o.packageName, "package", "", "Organization")
	newCmd.Flags().StringVar(&o.installPath, "installPath", "", "Install Folder")
	newCmd.Flags().StringVar(&o.token, "token", "", "Github Token")

	_ = newCmd.MarkFlagRequired("package")
	_ = newCmd.MarkFlagRequired("installPath")

	return newCmd
}
