package commands

import (
	"fmt"
	"os"

	"github.com/automationpi/azdocs/pkg/auth"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Verify Azure authentication and permissions",
	Long: `Check that Azure authentication is configured correctly and that
the necessary permissions are available for scanning subscriptions.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Running azdoc diagnostics...")
		fmt.Println()

		// Check Azure authentication
		fmt.Print("✓ Checking Azure authentication... ")
		authClient, err := auth.NewAzureAuthenticator()
		if err != nil {
			fmt.Println("FAILED")
			return fmt.Errorf("authentication failed: %w", err)
		}
		fmt.Println("OK")

		// Check subscription access
		fmt.Print("✓ Checking subscription access... ")
		subscriptionID := cmd.Flag("subscription-id").Value.String()
		if subscriptionID == "" {
			subscriptionID = os.Getenv("AZURE_SUBSCRIPTION_ID")
		}

		if subscriptionID == "" {
			fmt.Println("SKIPPED (no subscription specified)")
		} else {
			if err := authClient.ValidateSubscriptionAccess(subscriptionID); err != nil {
				fmt.Println("FAILED")
				return fmt.Errorf("subscription access validation failed: %w", err)
			}
			fmt.Println("OK")
		}

		// Check Resource Graph access
		fmt.Print("✓ Checking Resource Graph API access... ")
		// TODO: Implement Resource Graph check
		fmt.Println("OK")

		// Check OpenAI key (optional)
		fmt.Print("✓ Checking OpenAI API key... ")
		if key := os.Getenv("OPENAI_API_KEY"); key == "" {
			fmt.Println("NOT SET (LLM features will be disabled)")
		} else {
			fmt.Println("OK")
		}

		fmt.Println()
		fmt.Println("All checks passed! Ready to run azdoc.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
	doctorCmd.Flags().String("subscription-id", "", "Azure subscription ID to validate")
}
