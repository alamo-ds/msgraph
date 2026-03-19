package cmd

import (
	"context"
	"io"

	"github.com/spf13/cobra"
)

var usersCmd = &cobra.Command{
	Use:   "users",
	Args:  cobra.ExactArgs(1),
	Short: "get user info",
}

var (
	userId       string
	userEmail    string
	selectParams []string
)

func init() {
	rootCmd.AddCommand(usersCmd)

	usersCmd.AddCommand(usersGetCmd)
	usersGetCmd.Flags().StringVar(&userId, "id", "", "Microsoft user ID")
	usersGetCmd.Flags().StringVar(&userEmail, "email", "", "Outlook address associated with the user")
	usersGetCmd.Flags().StringArrayVar(&selectParams, "select", nil, "comma-separated values of field names to include in request")
	usersGetCmd.MarkFlagsMutuallyExclusive("id", "email")
}

var usersGetCmd = &cobra.Command{
	Use:   "get",
	Short: "get user info",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		out := cmd.OutOrStdout()

		switch {
		default:
			return handleGetUsers(ctx, out)
		case userId != "":
			return handleGetUserById(ctx, out)
		}
	},
}

func handleGetUsers(ctx context.Context, w io.Writer) error {
	users, err := client.Users().Select(selectParams...).Get(ctx)
	if err != nil {
		return err
	}

	jsonPrint(w, users)
	return nil
}

func handleGetUserById(ctx context.Context, w io.Writer) error {
	user, err := client.Users().Select(selectParams...).ById(userId).Get(ctx)
	if err != nil {
		return err
	}

	jsonPrint(w, user)
	return nil
}
