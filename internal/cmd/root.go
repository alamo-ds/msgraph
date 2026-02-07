package cmd

import (
	"context"
	"io"
	"os/exec"

	"github.com/alamo-ds/msgraph/graph"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:                "msgraph",
	Short:              "interact with your Microsoft Graph resources",
	PersistentPreRunE:  clientPreRun,
	PersistentPostRunE: clientPostRun,
}

func Execute(args []string, in io.Reader, out, err io.Writer) int {
	ctx := context.Background()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode()
		} else {
			return 1
		}
	}

	return 0
}

// TODO: do real pre-run checks
func clientPreRun(cmd *cobra.Command, args []string) error {
	client = graph.NewClient()
	return nil
}

func clientPostRun(cmd *cobra.Command, args []string) error {
	client.Close()
	return nil
}
