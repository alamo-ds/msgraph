package cmd

import "github.com/spf13/cobra"

var postsCmd = &cobra.Command{
	Use:     "posts",
	PreRunE: cobra.NoArgs,
	Short:   "get posts for a group ID & thread ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		posts, err := client.Groups().ById(groupId).Threads().ById(threadId).Get(cmd.Context())
		if err != nil {
			return err
		}

		jsonPrint(cmd.OutOrStdout(), posts)
		return nil
	},
}

var (
	threadId string
)

func init() {
	rootCmd.AddCommand(postsCmd)

	postsCmd.Flags().StringVar(&groupId, "group-id", "", "ID of the group")
	postsCmd.MarkFlagRequired("group-id")
	postsCmd.Flags().StringVar(&threadId, "thread-id", "", "thread whose posts to get")
	postsCmd.MarkFlagRequired("thread-id")
}
