package commands

import (
	"fmt"

	"github.com/automationpi/azdocs/pkg/graph"
	"github.com/automationpi/azdocs/pkg/renderer"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build documentation from cached data",
	Long: `Generate Markdown documentation and network diagrams from previously
scanned Azure resource data. Does not require network access.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Parse flags
		inDir := cmd.Flag("in").Value.String()
		outDir := cmd.Flag("out").Value.String()
		withDiagrams, _ := cmd.Flags().GetBool("with-diagrams")
		mdName := cmd.Flag("md-name").Value.String()
		theme := cmd.Flag("theme").Value.String()
		enableAI, _ := cmd.Flags().GetBool("enable-ai")
		openaiKey := cmd.Flag("openai-key").Value.String()

		fmt.Printf("Building documentation from %s...\n", inDir)

		// Load normalized data
		data, err := graph.LoadNormalizedData(inDir)
		if err != nil {
			return fmt.Errorf("failed to load data: %w", err)
		}

		// Build graph
		fmt.Println("Building topology graph...")
		graphBuilder := graph.NewBuilder(data)
		topology, err := graphBuilder.Build()
		if err != nil {
			return fmt.Errorf("failed to build graph: %w", err)
		}

		// Save graph
		if err := topology.SaveToFile(inDir + "/graph.json"); err != nil {
			return fmt.Errorf("failed to save graph: %w", err)
		}

		// Render documentation
		fmt.Println("Generating Markdown documentation...")
		mdRenderer := renderer.NewMarkdownRenderer(renderer.Config{
			OutputDir:    outDir,
			FileName:     mdName,
			WithDiagrams: withDiagrams,
			Theme:        theme,
			EnableAI:     enableAI,
			OpenAIKey:    openaiKey,
		})

		if err := mdRenderer.Render(topology); err != nil {
			return fmt.Errorf("failed to render markdown: %w", err)
		}

		// Generate diagrams if requested
		if withDiagrams {
			fmt.Println("Generating Draw.io diagrams...")
			if enableAI {
				fmt.Println("🤖 AI-enhanced diagram generation enabled")
			}
			diagramRenderer := renderer.NewDiagramRenderer(renderer.DiagramConfig{
				OutputDir: outDir + "/diagrams",
				Theme:     theme,
				EnableAI:  enableAI,
				OpenAIKey: openaiKey,
			})

			if err := diagramRenderer.RenderAll(topology); err != nil {
				return fmt.Errorf("failed to render diagrams: %w", err)
			}
		}

		fmt.Printf("\nBuild complete!\n")
		fmt.Printf("  Documentation: %s/%s\n", outDir, mdName)
		if withDiagrams {
			fmt.Printf("  Diagrams: %s/diagrams/\n", outDir)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().String("in", "./data", "input directory with cached JSON data")
	buildCmd.Flags().String("out", "./docs", "output directory for documentation")
	buildCmd.Flags().Bool("with-diagrams", true, "generate Draw.io diagrams")
	buildCmd.Flags().String("md-name", "SUBSCRIPTION.md", "output Markdown filename")
	buildCmd.Flags().String("theme", "default", "diagram theme")
	buildCmd.Flags().Bool("enable-ai", false, "enable AI-powered diagram optimization (requires OpenAI API key)")
	buildCmd.Flags().String("openai-key", "", "OpenAI API key (or set OPENAI_API_KEY env var)")
}
