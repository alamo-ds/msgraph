package cmd

import (
	"context"
	"io"

	"github.com/alamo-ds/msgraph/graph"
	"github.com/spf13/cobra"
)

var client *graph.Client

var groupsCmd = &cobra.Command{
	Use:   "groups",
	Args:  cobra.ExactArgs(1),
	Short: "interact with Groups and Group resources",
}

var (
	groupId    string
	getThreads bool
)

func init() {
	rootCmd.AddCommand(groupsCmd)

	groupsCmd.AddCommand(groupsGetCmd)

	groupsGetCmd.Flags().StringVar(&groupId, "id", "", "Microsoft Group ID")
	groupsGetCmd.Flags().BoolVar(&getThreads, "threads", false, "get threads associated with the group ID")
}

var groupsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "get info about group(s)",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		out := cmd.OutOrStdout()

		switch {
		default:
			return handleGetGroups(ctx, out)
		case groupId != "":
			if getThreads {
				return handleGetGroupThreads(ctx, out)
			}
			return handleGetGroupbyId(ctx, out)
		}
	},
}

func handleGetGroupbyId(ctx context.Context, w io.Writer) error {
	group, err := client.Groups().ById(groupId).Get(ctx)
	if err != nil {
		return err
	}

	jsonPrint(w, group)
	return nil
}

func handleGetGroups(ctx context.Context, w io.Writer) error {
	groups, err := client.Groups().Get(ctx)
	if err != nil {
		return err
	}

	jsonPrint(w, groups)
	return nil
}

func handleGetGroupThreads(ctx context.Context, w io.Writer) error {
	threads, err := client.Groups().ById(groupId).Threads().Get(ctx)
	if err != nil {
		return err
	}

	jsonPrint(w, threads)
	return nil
}
