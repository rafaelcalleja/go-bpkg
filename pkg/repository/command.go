package repository

import (
	"github.com/spf13/cobra"
)

type CommandProvider struct {
	Command      *cobra.Command
	preExecuteFn preExecute
}

type preExecute func(command *cobra.Command, releaseVersion *ReleaseVersion, tempDirectory string)

func NewCommandProvider(command *cobra.Command, preExecute preExecute) *CommandProvider {
	return &CommandProvider{
		command,
		preExecute,
	}
}

func (provider *CommandProvider) Download(releaseVersion *ReleaseVersion, tempDirectory string) error {
	cmdDownload := provider.Command

	if nil != provider.preExecuteFn {
		provider.preExecuteFn(cmdDownload, releaseVersion, tempDirectory)
	}

	return cmdDownload.Execute()
}
