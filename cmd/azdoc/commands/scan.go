package commands

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/automationpi/azdocs/pkg/auth"
	"github.com/automationpi/azdocs/pkg/cache"
	"github.com/automationpi/azdocs/pkg/discovery"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan Azure subscription and cache resource data",
	Long: `Connect to Azure, discover all resources, and cache the raw data locally.
This command performs inventory collection via Azure Resource Graph and
detailed queries via ARM APIs.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Parse flags
		subscriptionID := cmd.Flag("subscription-id").Value.String()
		if subscriptionID == "" {
			subscriptionID = os.Getenv("AZURE_SUBSCRIPTION_ID")
			if subscriptionID == "" {
				return fmt.Errorf("subscription ID is required (use --subscription-id or AZURE_SUBSCRIPTION_ID)")
			}
		}

		tenantID := cmd.Flag("tenant-id").Value.String()
		concurrency, _ := cmd.Flags().GetInt("concurrency")
		timeout, _ := cmd.Flags().GetDuration("timeout")
		cacheDir := cmd.Flag("cache-dir").Value.String()
		jsonOut := cmd.Flag("json-out").Value.String()
		noProgress, _ := cmd.Flags().GetBool("no-progress")

		// Initialize auth
		authClient, err := auth.NewAzureAuthenticator()
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		// Initialize cache
		cacheClient, err := cache.NewCache(cacheDir)
		if err != nil {
			return fmt.Errorf("cache initialization failed: %w", err)
		}
		defer cacheClient.Close()

		// Initialize discovery
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		discoveryClient := discovery.NewDiscoveryClient(authClient, cacheClient, discovery.Config{
			SubscriptionID: subscriptionID,
			TenantID:       tenantID,
			Concurrency:    concurrency,
			ShowProgress:   !noProgress,
		})

		// Run discovery
		if !noProgress {
			fmt.Printf("Scanning subscription %s...\n", subscriptionID)
		}

		result, err := discoveryClient.Discover(ctx)
		if err != nil {
			return fmt.Errorf("discovery failed: %w", err)
		}

		// Save results
		if err := result.SaveToDirectory(jsonOut); err != nil {
			return fmt.Errorf("failed to save results: %w", err)
		}

		// Fetch Azure Advisor recommendations
		if !noProgress {
			fmt.Println("\nFetching Azure Advisor recommendations...")
		}
		recommendations, err := discoveryClient.FetchRecommendations(ctx)
		if err != nil {
			fmt.Printf("⚠️  Warning: Failed to fetch recommendations: %v\n", err)
			fmt.Println("   Continuing without recommendations...")
		} else {
			// Save recommendations
			recPath := jsonOut + "/raw/recommendations.json"
			if err := discovery.SaveRecommendations(recommendations, recPath); err != nil {
				fmt.Printf("⚠️  Warning: Failed to save recommendations: %v\n", err)
			} else if !noProgress {
				fmt.Printf("  ✅ Found %d recommendations\n", len(recommendations))
			}
		}

		if !noProgress {
			fmt.Printf("\nScan complete!\n")
			fmt.Printf("  Resources discovered: %d\n", result.Stats.TotalResources)
			fmt.Printf("  VNets: %d\n", result.Stats.VNets)
			fmt.Printf("  Subnets: %d\n", result.Stats.Subnets)
			fmt.Printf("  NSGs: %d\n", result.Stats.NSGs)
			fmt.Printf("  Route Tables: %d\n", result.Stats.RouteTables)
			fmt.Printf("\nData saved to: %s\n", jsonOut)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)

	scanCmd.Flags().String("subscription-id", "", "Azure subscription ID to scan")
	scanCmd.Flags().String("tenant-id", "", "Azure tenant ID (optional)")
	scanCmd.Flags().Int("concurrency", 8, "number of concurrent API requests")
	scanCmd.Flags().Duration("timeout", 5*time.Minute, "overall operation timeout")
	scanCmd.Flags().String("cache-dir", ".azdoc", "cache directory path")
	scanCmd.Flags().String("json-out", "./data", "output directory for JSON files")
	scanCmd.Flags().Bool("no-progress", false, "suppress progress indicators")

	// Bind to viper
	viper.BindPFlag("subscription-id", scanCmd.Flags().Lookup("subscription-id"))
	viper.BindPFlag("concurrency", scanCmd.Flags().Lookup("concurrency"))
	viper.BindPFlag("cache-dir", scanCmd.Flags().Lookup("cache-dir"))
}
