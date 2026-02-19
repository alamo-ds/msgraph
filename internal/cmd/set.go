package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/alamo-ds/msgraph/env"
	"github.com/alamo-ds/msgraph/graph"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "set tenant ID and client ID/secret from env or manual entry",
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		var err error

		if tenantId == "" {
			err = flagErr("tenant-id")
		} else if clientId == "" {
			err = flagErr("client-id")
		} else if clientSecret == "" {
			err = flagErr("client-secret")
		}

		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := graph.AzureADConfig{
			TenantID:     tenantId,
			ClientID:     clientId,
			ClientSecret: clientSecret,
			Scopes:       []string{graph.DefaultScopes},
		}

		data, _ := json.MarshalIndent(cfg, "", "  ")
		env.WriteConfigFile(data)
		fmt.Fprintf(cmd.OutOrStdout(), "config for tenant ID %s stored successfully\n", tenantId)
	},
}

var (
	tenantId     string
	clientId     string
	clientSecret string
)

func init() {
	rootCmd.AddCommand(setCmd)

	setCmd.Flags().StringVar(&tenantId, "tenant-id", "", "")
	setCmd.Flags().StringVar(&clientId, "client-id", "", "")
	setCmd.Flags().StringVar(&clientSecret, "client-secret", "", "")

	tenantId = os.Getenv("TENANT_ID")
	clientId = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")
}

func flagErr(envVar string) error {
	return fmt.Errorf("value for %s not provided", envVar)
}
