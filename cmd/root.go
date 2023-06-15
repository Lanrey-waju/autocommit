package cmd

import (
	"fmt"
	"os"

	"github.com/christian-gama/autocommit/internal/autocommit"
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "",
	Run:   runCmd,
	Short: "Autocommit is a CLI tool that uses OpenAI's models to generate commit messages based on the changes made in the repository.",
}

// Execute executes the root command.
func Execute() error {
	return cmd.Execute()
}

func runCmd(cmd *cobra.Command, args []string) {
	err := verifyConfigCommand.Execute(askConfigsCli.Execute)
	if err != nil {
		panic(err)
	}

	response, err := generatorCommand.Execute()
	if err != nil {
		panic(err)
	}

	fmt.Printf("📝 Commit message generated: \n%s\n\n", response)

	option, err := postCommitCli.Execute()
	if err != nil {
		panic(err)
	}

	switch option {
	case autocommit.CommitChangesOption:
		if err := commitCommand.Execute(response); err != nil {
			panic(err)
		}

	case autocommit.CopyToClipboardOption:
		if err := clipboardCommand.Execute(response); err != nil {
			panic(err)
		}

	case autocommit.RegenerateOption:
		runCmd(cmd, args)

	case autocommit.ExitOption:
		os.Exit(0)
	}
}
