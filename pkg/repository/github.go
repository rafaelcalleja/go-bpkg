package repository

import (
	"fmt"
	"github.com/cli/cli/v2/pkg/cmd/release/download"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
)

const GithubRepository = "github.com"

type GithubProvider struct {
	*CommandProvider
	factory  *cmdutil.Factory
	hostname string
	archive  string
}

func NewGithubProviderWith(options ...func(*GithubProvider) error) (*GithubProvider, error) {
	var githubProvider = new(GithubProvider)

	for _, option := range options {
		err := option(githubProvider)
		if err != nil {
			return nil, err
		}
	}

	if nil == githubProvider.CommandProvider {
		cmdDownload := download.NewCmdDownload(githubProvider.factory, nil)
		cmdDownload.PersistentFlags().StringP("repo", "R", "", "")

		githubProvider.CommandProvider = NewCommandProvider(
			cmdDownload,
			createGithubProviderPreExecution(githubProvider),
		)
	}

	return githubProvider, nil
}

func WithCommandProvider(provider *CommandProvider) func(*GithubProvider) error {
	return func(g *GithubProvider) error {
		g.CommandProvider = provider
		return nil
	}
}

func WithFactory(factory *cmdutil.Factory) func(*GithubProvider) error {
	return func(g *GithubProvider) error {
		g.factory = factory
		return nil
	}
}

func WithHostname(hostname string) func(*GithubProvider) error {
	return func(g *GithubProvider) error {
		g.hostname = hostname
		return nil
	}
}

func WithArchiveFormat(archive string) func(*GithubProvider) error {
	return func(g *GithubProvider) error {
		g.archive = archive
		return nil
	}
}

func NewGithubProvider(factory *cmdutil.Factory) *GithubProvider {
	provider, _ := NewGithubProviderWith(
		WithHostname(GithubRepository),
		WithFactory(factory),
		WithArchiveFormat("tar.gz"),
	)

	return provider
}

func createGithubProviderPreExecution(provider *GithubProvider) preExecute {
	return func(command *cobra.Command, releaseVersion *ReleaseVersion, tempDirectory string) {
		command.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
			provider.factory.BaseRepo = cmdutil.OverrideBaseRepoFunc(provider.factory, fmt.Sprintf("%s/%s/%s", provider.hostname, releaseVersion.Organization, releaseVersion.Name))
			return nil
		}

		args := []string{
			releaseVersion.Version(),
			"--dir",
			tempDirectory,
			"--archive",
			provider.archive,
		}

		command.SetArgs(args)
	}
}
