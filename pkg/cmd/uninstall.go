package cmd

import (
	"fmt"
	"github.com/rafaelcalleja/go-bpkg/pkg/repository"
	"github.com/rafaelcalleja/go-kit/cmd/helper"
	"github.com/rafaelcalleja/go-kit/cmd/termcolor"
	"github.com/rafaelcalleja/go-kit/logger"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

func NewPackageUninstall(
	helper helper.ErrorHelper,
	log logger.Logger,
	term termcolor.TermColor,
) *cobra.Command {
	o := &PackageInstallOptions{}

	newCmd := &cobra.Command{
		Use:   "uninstall",
		Short: "BPKG uninstall",
		Run: func(cmd *cobra.Command, args []string) {
			if "" != strings.TrimSpace(o.token) {
				_ = os.Setenv("GITHUB_TOKEN", o.token)
			}

			fqpVO, err := repository.NewFullyQualifyPackage(o.packageName)
			helper.CheckErr(err)
			pkgName := fmt.Sprintf("%s-%s", fqpVO.Organization(), fqpVO.Name())

			packagesInstalled, err := repository.PackagesInstalled(o.installPath)
			helper.CheckErr(err)

			for _, pkg := range packagesInstalled {
				if pkg.Name == pkgName {
					err = pkg.Uninstall(filepath.Join(o.installPath, pkgName))
					helper.CheckErr(err)
					log.Infof("Package %s:%s uninstalled!", term.ColorInfo(fmt.Sprintf("%s/%s", fqpVO.Organization(), fqpVO.Name())),
						term.ColorInfo(fqpVO.Version()))

					break
				}
			}
		},
	}

	newCmd.Flags().StringVar(&o.packageName, "package", "", "[package to uninstall] package/name:v1.0.0")
	newCmd.Flags().StringVar(&o.installPath, "installPath", "./deps", "[package install path]")

	_ = newCmd.MarkFlagRequired("package")

	return newCmd
}
