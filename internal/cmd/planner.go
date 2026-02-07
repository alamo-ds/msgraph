package cmd

import (
	"context"
	"errors"
	"io"

	"github.com/alamo-ds/msgraph/graph"
	"github.com/spf13/cobra"
)

var plannerCmd = &cobra.Command{
	Use:   "planner",
	Args:  cobra.ExactArgs(1),
	Short: "interact with Planner resources",
}

var (
	plannerId          string
	updatePlannerTitle string
	getTasks           bool
	getBuckets         bool
)

func init() {
	rootCmd.AddCommand(plannerCmd)

	plannerCmd.PersistentFlags().StringVar(&plannerId, "id", "", "ID of the Planner plan")
	plannerCmd.PersistentFlags().StringVar(&groupId, "group-id", "", "Microsoft Group ID for which to fetch Planner plans")

	plannerCmd.AddCommand(plannerGetCmd, plannerUpdateCmd)

	plannerGetCmd.Flags().BoolVar(&getTasks, "tasks", false, "Return tasks item for plan ID")
	plannerGetCmd.Flags().BoolVar(&getBuckets, "buckets", false, "Return all buckets for plan ID")
	plannerGetCmd.MarkFlagsOneRequired("id", "group-id")
	plannerGetCmd.MarkFlagsMutuallyExclusive("id", "group-id")
	plannerGetCmd.MarkFlagsMutuallyExclusive("tasks", "buckets")

	plannerUpdateCmd.MarkPersistentFlagRequired("id")
	plannerUpdateCmd.Flags().StringVar(&updatePlannerTitle, "title", "", "new title to assign to Planner plan")
}

var plannerGetCmd = &cobra.Command{
	Use:   "get",
	Short: "get info about Planner plan(s)",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		out := cmd.OutOrStdout()

		switch {
		default:
			panic("neither id nor group-id flag is set")
		case plannerId != "":
			if getTasks {
				return handleGetTasksForPlan(ctx, out)
			}
			if getBuckets {
				return handleGetBucketsForPlan(ctx, out)
			}

			return handleGetPlan(ctx, out)
		case groupId != "":
			return handleGetPlansForGroup(ctx, out)
		}
	},
}

func handleGetPlansForGroup(ctx context.Context, w io.Writer) error {
	plans, err := client.Groups().ById(groupId).Plans().Get(ctx)
	if err != nil {
		return err
	}

	jsonPrint(w, plans)
	return nil
}

func handleGetPlan(ctx context.Context, w io.Writer) error {
	plan, err := client.Planner().ById(plannerId).Get(ctx)
	if err != nil {
		return err
	}

	jsonPrint(w, plan)
	return nil
}

func handleGetTasksForPlan(ctx context.Context, w io.Writer) error {
	tasks, err := client.Planner().ById(plannerId).Tasks().Get(ctx)
	if err != nil {
		return err
	}

	jsonPrint(w, tasks)
	return nil
}

func handleGetBucketsForPlan(ctx context.Context, w io.Writer) error {
	buckets, err := client.Planner().ById(plannerId).Buckets().Get(ctx)
	if err != nil {
		return err
	}

	jsonPrint(w, buckets)
	return nil
}

var plannerUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "update a Planner plan by id",
	RunE: func(cmd *cobra.Command, args []string) error {
		if updatePlannerTitle == "" {
			return errors.New("value for new title not set")
		}

		plan, err := client.Planner().ById(plannerId).Patch(cmd.Context(), graph.PatchPlanParams{
			Title: updatePlannerTitle,
		})
		if err != nil {
			return err
		}

		jsonPrint(cmd.OutOrStdout(), plan)
		return nil
	},
}
