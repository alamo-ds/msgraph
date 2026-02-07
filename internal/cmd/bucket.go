package cmd

import (
	"errors"
	"fmt"

	"github.com/alamo-ds/msgraph/graph"
	"github.com/spf13/cobra"
)

var bucketCmd = &cobra.Command{
	Use:   "bucket",
	Args:  cobra.ExactArgs(1),
	Short: "interact with Planner buckets",
}

var (
	bucketId         string
	updateBucketName string
)

func init() {
	rootCmd.AddCommand(bucketCmd)

	bucketCmd.AddCommand(bucketUpdateCmd)
	bucketCmd.MarkPersistentFlagRequired("id")

	bucketUpdateCmd.Flags().StringVar(&updateBucketName, "name", "", "new name to assign to Planner bucket")
	bucketUpdateCmd.Flags().StringVar(&bucketId, "id", "", "ID of the Planner bucket")
	bucketUpdateCmd.MarkFlagRequired("id")

	bucketCmd.AddCommand(bucketCreateCmd)
	bucketCreateCmd.Flags().StringVar(&plannerId, "plan-id", "", "Planner plan ID for which to create new bucket")
	bucketCreateCmd.MarkFlagRequired("plan-id")
	bucketCreateCmd.Flags().StringVar(&updateBucketName, "name", "", "new name to assign to Planner bucket")
}

var bucketUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "update a Planner bucket by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		if updateBucketName == "" {
			return errors.New("value for new name not set")
		}

		bucket, err := client.Planner().Buckets().ById(bucketId).Patch(cmd.Context(), graph.PatchBucketParams{
			Name: updateBucketName,
		})
		if err != nil {
			return err
		}

		jsonPrint(cmd.OutOrStdout(), bucket)
		return err
	},
}

var bucketCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create a new Planner bucket in the provided plan ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		if updateBucketName == "" {
			return errors.New("value for new name not set")
		}

		bucket, err := client.Planner().Buckets().Post(cmd.Context(), graph.PostBucketParams{
			Name:   updateBucketName,
			PlanID: plannerId,
		})
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		fmt.Fprintln(out, "new bucket created")
		jsonPrint(out, bucket)
		return nil
	},
}
