package cmd

import (
	"context"
	"io"
	"os"
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

var (
	tenantId     string
	clientId     string
	clientSecret string
)

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
	tenantId = os.Getenv("TENANT_ID")
	clientId = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")

	client = graph.NewClient(cmd.Context(), clientSecret)
	return nil
}

func clientPostRun(cmd *cobra.Command, args []string) error {
	client.Close()
	return nil
}
