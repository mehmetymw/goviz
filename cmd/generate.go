package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"goviz/pkg/graph"
	"goviz/pkg/output"
	"goviz/pkg/parser"

	"github.com/spf13/cobra"
)

var (
	format     string
	outputFile string
)

var generateCmd = &cobra.Command{
	Use:   "generate [path]",
	Short: "Generate dependency graph from go.mod file",
	Long: `Generate a dependency graph from a go.mod file.
	
If no path is provided, the current directory will be used.
The tool will look for go.mod file in the specified directory.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var projectPath string

		if len(args) == 0 {
			projectPath = "."
		} else {
			projectPath = args[0]
		}

		absPath, err := filepath.Abs(projectPath)
		if err != nil {
			return fmt.Errorf("failed to get absolute path: %w", err)
		}

		goModPath := filepath.Join(absPath, "go.mod")
		if _, err := os.Stat(goModPath); os.IsNotExist(err) {
			return fmt.Errorf("go.mod file not found in %s", absPath)
		}

		fmt.Printf("Parsing go.mod from %s...\n", absPath)
		modFile, err := parser.ParseGoMod(goModPath)
		if err != nil {
			return fmt.Errorf("failed to parse go.mod: %w", err)
		}

		goSumPath := filepath.Join(absPath, "go.sum")
		enhancedGraph, err := graph.BuildEnhancedDependencyGraph(modFile, goSumPath)
		if err != nil {
			return fmt.Errorf("failed to build enhanced dependency graph: %w", err)
		}

		enhancedGraph.DetectVersionConflicts()
		if err := enhancedGraph.AnalyzeLicenses(); err != nil {
			return fmt.Errorf("failed to analyze licenses: %w", err)
		}
		if err := enhancedGraph.CheckSecurity(); err != nil {
			return fmt.Errorf("failed to check security: %w", err)
		}

		switch format {
		case "dot":
			if outputFile == "" {
				outputFile = "depgraph.dot"
			}
			return output.GenerateEnhancedDOT(enhancedGraph, outputFile)
		case "png":
			if outputFile == "" {
				outputFile = "depgraph.png"
			}
			return output.GeneratePNG(enhancedGraph, outputFile)
		case "svg":
			if outputFile == "" {
				outputFile = "depgraph.svg"
			}
			return output.GenerateSVG(enhancedGraph, outputFile)
		case "json":
			return output.GenerateJSON(enhancedGraph, outputFile, absPath)
		case "yaml":
			return output.GenerateYAML(enhancedGraph, outputFile, absPath)
		case "tree", "ascii":
			return output.GenerateASCIITree(enhancedGraph.DependencyGraph)
		default:
			return fmt.Errorf("unsupported format: %s. Supported formats: dot, png, svg, json, yaml, tree, ascii", format)
		}
	},
}

func init() {
	generateCmd.Flags().StringVarP(&format, "format", "f", "tree", "Output format (dot, png, svg, json, yaml, tree, ascii)")
	generateCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file")
}
