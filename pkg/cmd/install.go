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
	packageName  string
	installPath  string
	token        string
	metadataJson string
	alias        string
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

			if "" == strings.TrimSpace(fqpVO.Version()) {
				log.Errorf("version is required, package format is [%s] || [%s]", term.ColorInfo("package/name:v1.0.0"), term.ColorInfo("package/name:latest"))

				return
			}

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

				fqpVO = fqpVO.CopyWithVersion(releaseVersion.Version())
			}

			if "" != o.alias {
				fqpVO = fqpVO.CopyWithName(o.alias)
			}

			pkgName := fmt.Sprintf("%s-%s", fqpVO.Organization(), fqpVO.Name())
			packagesInstalled, err := repository.PackagesInstalled(o.installPath)
			helper.CheckErr(err)

			for _, pkg := range packagesInstalled {
				if pkg.Name == pkgName && pkg.Version == fqpVO.Version() {
					log.Infof("Package %s already at version %s", term.ColorInfo(fqpVO.String()),
						term.ColorInfo(releaseVersion.Version()))

					return
				}
			}

			log.Infof("Installing Package %s at %s", term.ColorInfo(releaseVersion.String()),
				term.ColorInfo(o.installPath))

			assetGithub := repository.NewGithubProvider(factory)
			asset, err := releaseVersion.DownloadAsset(assetGithub, o.installPath)
			if "" != o.alias {
				asset = asset.CopyWithName(fmt.Sprintf("%s-%s", fqpVO.Organization(), o.alias))
			}

			helper.CheckErr(err)

			if "" != strings.TrimSpace(o.metadataJson) {
				metadata, err := repository.NewPackageInstallerFromLiteral(o.metadataJson)
				helper.CheckErr(err)

				err = asset.Install(metadata, o.installPath)
				helper.CheckErr(err)
			} else {
				err = releaseVersion.InstallAsset(asset, o.installPath)
				helper.CheckErr(err)
			}

			log.Infof("Installed Successfully")
		},
	}

	newCmd.Flags().StringVar(&o.packageName, "package", "", "[package to install] package/name:v1.0.0")
	newCmd.Flags().StringVar(&o.installPath, "installPath", "./deps", "[package install path]")
	newCmd.Flags().StringVar(&o.token, "token", "", "Github Token")
	newCmd.Flags().StringVar(&o.metadataJson, "metadataJson", "", "overwrite current package.json")
	newCmd.Flags().StringVar(&o.alias, "alias", "", "package name is replace using alias")

	_ = newCmd.MarkFlagRequired("package")

	return newCmd
}
