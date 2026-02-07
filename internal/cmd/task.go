package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/alamo-ds/msgraph/env"
	"github.com/alamo-ds/msgraph/graph"
	"github.com/spf13/cobra"
)

var taskCmd = &cobra.Command{
	Use:   "task",
	Args:  cobra.ExactArgs(1),
	Short: "interact with Planner tasks",
}

var (
	taskId          string
	details         bool
	updateTaskTitle string
	taskFile        string
)

func init() {
	rootCmd.AddCommand(taskCmd)

	taskCmd.AddCommand(taskGetCmd)
	taskGetCmd.PersistentFlags().StringVar(&taskId, "id", "", "ID of the Planner task")
	taskCmd.Flags().BoolVar(&details, "details", false, "fetch details about the Planner task")
	taskCmd.MarkPersistentFlagRequired("id")

	taskCmd.AddCommand(taskUpdateCmd)
	taskUpdateCmd.Flags().StringVar(&updateTaskTitle, "title", "", "new title to assign to Planner task")

	taskCmd.AddCommand(taskCreateCmd)
	taskCreateCmd.Flags().StringVar(&plannerId, "plan-id", "", "plan ID to which to add task")
	taskCreateCmd.Flags().StringVarP(&taskFile, "file", "f", "", "file from which to read new tasks")
	taskCreateCmd.MarkFlagRequired("plan-id")
}

var taskGetCmd = &cobra.Command{
	Use:   "get",
	Short: "get a Planner task by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		out := cmd.OutOrStdout()

		if details {
			taskDetails, err := client.Planner().Tasks().ById(taskId).Details().Get(ctx)
			if err != nil {
				return err
			}

			jsonPrint(out, taskDetails)
		} else {
			task, err := client.Planner().Tasks().ById(taskId).Get(ctx)
			if err != nil {
				return err
			}

			jsonPrint(out, task)
		}

		return nil
	},
}

var taskUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "update a Planner task by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: add a config flag, obviously there are more things
		// than just the title that we can update.
		if updateTaskTitle == "" {
			return errors.New("value for new title not set")
		}

		task, err := client.Planner().Tasks().ById(taskId).Patch(cmd.Context(), graph.PatchTaskParams{
			Title: updateTaskTitle,
		})
		if err != nil {
			return err
		}

		jsonPrint(cmd.OutOrStdout(), task)
		return nil
	},
}

var taskCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create a new Planner task for the provided plan ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		var in = cmd.InOrStdin()

		if taskFile != "" {
			f, err := env.SafeOpen(taskFile)
			if err != nil {
				return fmt.Errorf("os.Open: %v", err)
			}

			in = f
		} else {
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) != 0 {
				return errors.New("no data provided. Either pipe JSON into program or specify file flag")
			}
		}

		var params graph.PostTaskParams
		if err := json.NewDecoder(in).Decode(&params); err != nil {
			return fmt.Errorf("couldn't decode create task params JSON: %v", err)
		}
		params.PlanID = plannerId

		ret, err := client.Planner().Tasks().Post(cmd.Context(), params)
		if err != nil {
			return err
		}

		jsonPrint(cmd.OutOrStdout(), ret)
		return nil
	},
}
