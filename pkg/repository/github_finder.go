package repository

import (
	"bytes"
	"fmt"
	"github.com/cli/cli/v2/pkg/cmd/release/list"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

type GithubVersionFinder struct {
	command *cobra.Command
	factory *cmdutil.Factory
	limit   int
}

func NewGithubVersionFinderWith(options ...func(*GithubVersionFinder) error) (*GithubVersionFinder, error) {
	var githubVersionFinder = new(GithubVersionFinder)

	for _, option := range options {
		err := option(githubVersionFinder)
		if err != nil {
			return nil, err
		}
	}

	if nil == githubVersionFinder.command {
		cmd := list.NewCmdList(githubVersionFinder.factory, nil)
		cmd.PersistentFlags().StringP("repo", "R", "", "")

		githubVersionFinder.command = cmd
	}

	return githubVersionFinder, nil
}

func (g *GithubVersionFinder) Latest(organization string, name string) (string, error) {
	stdout := g.factory.IOStreams.Out
	buf := new(bytes.Buffer)

	g.command.PersistentPreRunE = func(rootCmd *cobra.Command, args []string) error {
		g.factory.IOStreams.SetColorEnabled(false)
		g.factory.BaseRepo = cmdutil.OverrideBaseRepoFunc(g.factory, fmt.Sprintf("%s/%s/%s", "github.com", organization, name))
		g.factory.IOStreams.Out = buf

		return nil
	}

	g.command.SetArgs([]string{"-L", strconv.Itoa(g.limit)})

	err := g.command.Execute()
	if nil != err {
		return "", nil
	}

	g.factory.IOStreams.Out = stdout

	return strings.Fields(buf.String())[0], nil
}

func FinderWithCommand(command *cobra.Command) func(*GithubVersionFinder) error {
	return func(g *GithubVersionFinder) error {
		g.command = command
		return nil
	}
}

func FinderWithFactory(factory *cmdutil.Factory) func(*GithubVersionFinder) error {
	return func(g *GithubVersionFinder) error {
		g.factory = factory
		return nil
	}
}

func FinderWithLimit(limit int) func(*GithubVersionFinder) error {
	return func(g *GithubVersionFinder) error {
		g.limit = limit
		return nil
	}
}

func NewGithubVersionFinder(factory *cmdutil.Factory) *GithubVersionFinder {
	finder, _ := NewGithubVersionFinderWith(
		FinderWithFactory(factory),
		FinderWithLimit(1),
	)

	return finder
}
