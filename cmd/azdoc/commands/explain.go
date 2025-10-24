package commands

import (
	"fmt"
	"os"

	"github.com/automationpi/azdocs/pkg/llm"
	"github.com/spf13/cobra"
)

var explainCmd = &cobra.Command{
	Use:   "explain",
	Short: "Add LLM-generated explanations to documentation",
	Long: `Use LLM to add narrative sections, risk analysis, and insights
to existing documentation. Requires OpenAI API key.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Parse flags
		inDir := cmd.Flag("in").Value.String()
		docPath := cmd.Flag("doc").Value.String()
		model := cmd.Flag("model").Value.String()
		maxTokens, _ := cmd.Flags().GetInt("max-tokens")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		// Check API key
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			return fmt.Errorf("OPENAI_API_KEY environment variable is required")
		}

		fmt.Printf("Generating LLM explanations for %s...\n", docPath)

		// Initialize LLM client
		llmClient := llm.NewClient(llm.Config{
			Enabled: true,
			APIKey:  apiKey,
			Model:   model,
		})

		// Note: maxTokens and dryRun parameters are currently not used
		_ = maxTokens
		_ = dryRun

		// Load normalized data
		data, err := loadDataForLLM(inDir)
		if err != nil {
			return fmt.Errorf("failed to load data: %w", err)
		}

		// TODO: Implement GenerateExplanations method in LLM client
		_ = llmClient
		_ = data
		_ = dryRun
		fmt.Println("⚠️  Explain command is currently being refactored for AI integration")
		fmt.Println("    Use 'build --enable-ai' for AI-enhanced diagrams")
		return nil

		// Generate explanations with 2-parse validation
		// fmt.Println("Running LLM with structured output validation...")
		// explanations, err := llmClient.GenerateExplanations(data)
		// if err != nil {
		// 	return fmt.Errorf("LLM generation failed: %w", err)
		// }

		// if dryRun {
		// 	fmt.Println("\n--- DRY RUN: Generated explanations ---")
		// 	fmt.Printf("Overview: %s\n", explanations.Overview)
		// 	fmt.Printf("Key Risks: %d items\n", len(explanations.Risks))
		// 	fmt.Printf("Questions: %d items\n", len(explanations.Questions))
		// 	return nil
		// }

		// Inject explanations into document
		// fmt.Println("Injecting explanations into documentation...")
		// if err := injectExplanations(docPath, explanations); err != nil {
		// 	return fmt.Errorf("failed to inject explanations: %w", err)
		// }

		// fmt.Println("\nExplanations added successfully!")
		// return nil
	},
}

func init() {
	rootCmd.AddCommand(explainCmd)

	explainCmd.Flags().String("in", "./data", "input directory with cached JSON data")
	explainCmd.Flags().String("doc", "./docs/SUBSCRIPTION.md", "documentation file to enhance")
	explainCmd.Flags().String("model", "gpt-4o-mini", "OpenAI model to use")
	explainCmd.Flags().Int("max-tokens", 2000, "maximum tokens for LLM response")
	explainCmd.Flags().Bool("dry-run", false, "preview explanations without modifying files")
}

// Helper functions (stubs for now)
func loadDataForLLM(inDir string) (interface{}, error) {
	// TODO: Implement
	return nil, nil
}

func injectExplanations(docPath string, explanations interface{}) error {
	// TODO: Implement
	return nil
}
