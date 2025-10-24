package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Run scan, build, and optionally explain in sequence",
	Long: `Convenience command that runs the complete pipeline:
  1. Scan Azure subscription
  2. Build documentation and diagrams
  3. Add LLM explanations (if --with-llm is set)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		withLLM, _ := cmd.Flags().GetBool("with-llm")

		fmt.Println("=== Phase 1: Scanning ===")
		if err := scanCmd.RunE(cmd, args); err != nil {
			return fmt.Errorf("scan failed: %w", err)
		}

		fmt.Println("\n=== Phase 2: Building ===")
		if err := buildCmd.RunE(cmd, args); err != nil {
			return fmt.Errorf("build failed: %w", err)
		}

		if withLLM {
			fmt.Println("\n=== Phase 3: LLM Explanations ===")
			if err := explainCmd.RunE(cmd, args); err != nil {
				return fmt.Errorf("explain failed: %w", err)
			}
		}

		fmt.Println("\n=== All phases complete! ===")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(allCmd)

	// Add all flags from subcommands directly
	allCmd.Flags().Bool("with-llm", false, "include LLM explanations")

	// Scan flags
	allCmd.Flags().String("subscription-id", "", "Azure subscription ID to scan")
	allCmd.Flags().String("tenant-id", "", "Azure tenant ID (optional)")
	allCmd.Flags().Int("concurrency", 8, "number of concurrent API requests")
	allCmd.Flags().Duration("timeout", 5*60*1000000000, "overall operation timeout (default 5m)")
	allCmd.Flags().String("cache-dir", ".azdoc", "cache directory path")
	allCmd.Flags().String("json-out", "./data", "output directory for JSON files")
	allCmd.Flags().Bool("no-progress", false, "suppress progress indicators")

	// Build flags
	allCmd.Flags().String("in", "./data", "input directory with cached JSON data")
	allCmd.Flags().String("out", "./docs", "output directory for documentation")
	allCmd.Flags().Bool("with-diagrams", true, "generate Draw.io diagrams")
	allCmd.Flags().String("md-name", "SUBSCRIPTION.md", "output Markdown filename")
	allCmd.Flags().String("theme", "default", "diagram theme")
	allCmd.Flags().Bool("enable-ai", false, "enable AI-powered diagram optimization (requires OpenAI API key)")
	allCmd.Flags().String("openai-key", "", "OpenAI API key (or set OPENAI_API_KEY env var)")

	// Explain flags
	allCmd.Flags().String("doc", "./docs/SUBSCRIPTION.md", "documentation file to enhance")
	allCmd.Flags().String("model", "gpt-4o-mini", "OpenAI model to use")
	allCmd.Flags().Int("max-tokens", 2000, "maximum tokens for LLM response")
	allCmd.Flags().Bool("dry-run", false, "preview explanations without modifying files")
}
