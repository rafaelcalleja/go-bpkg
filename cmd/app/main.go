//nolint:gofmt,goimports
//go:build !windows
// +build !windows

package app

import (
	"github.com/rafaelcalleja/go-bpkg/pkg/cmd"
)

// Run runs the command, if args are not nil they will be set on the command
func Run(args []string) error {
	rootCmd := cmd.Main()
	if args != nil {
		args = args[1:]
		rootCmd.SetArgs(args)
	}

	return rootCmd.Execute()
}
