package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"goviz/pkg/graph"
	"goviz/pkg/output"
	"goviz/pkg/parser"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	analyzeFormat string
	analyzeOutput string
	showConflicts bool
	showOutdated  bool
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze [path]",
	Short: "Analyze dependencies for conflicts, security issues, and health",
	Long: `Perform comprehensive analysis of your Go module dependencies.
	
This command analyzes:
- Version conflicts
- Security vulnerabilities  
- License compatibility
- Outdated packages
- Dependency health metrics`,
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

		fmt.Printf("Analyzing dependencies from %s...\n", absPath)
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

		switch analyzeFormat {
		case "json":
			return output.GenerateJSON(enhancedGraph, analyzeOutput, absPath)
		case "yaml":
			return output.GenerateYAML(enhancedGraph, analyzeOutput, absPath)
		case "text", "console":
			return generateAnalysisReport(enhancedGraph)
		default:
			return fmt.Errorf("unsupported format: %s. Supported formats: json, yaml, text, console", analyzeFormat)
		}
	},
}

func generateAnalysisReport(graph *graph.EnhancedDependencyGraph) error {

	red := color.New(color.FgRed, color.Bold)
	green := color.New(color.FgGreen, color.Bold)
	yellow := color.New(color.FgYellow, color.Bold)
	blue := color.New(color.FgBlue, color.Bold)

	blue.Printf("🔍 Dependency Analysis Report\n")
	blue.Printf("============================\n\n")

	fmt.Printf("Module: %s\n", graph.ModuleName)
	if graph.ModuleGoVersion != "" {
		fmt.Printf("Go Version: %s\n", graph.ModuleGoVersion)
	}
	fmt.Println()

	stats := graph.GetStatistics()
	blue.Printf("📊 Statistics:\n")
	fmt.Printf("  Total Dependencies: %v\n", stats["total_dependencies"])
	fmt.Printf("  Direct Dependencies: %v\n", stats["direct_dependencies"])
	fmt.Printf("  Indirect Dependencies: %v\n", stats["indirect_dependencies"])
	fmt.Printf("  Transitive Dependencies: %v\n", stats["transitive_dependencies"])
	fmt.Printf("  Unique Licenses: %v\n", stats["unique_licenses"])
	fmt.Println()

	if len(graph.Conflicts) > 0 {
		red.Printf("⚡ Version Conflicts (%d):\n", len(graph.Conflicts))
		for _, conflict := range graph.Conflicts {
			fmt.Printf("  • %s: %s vs %s (%s)\n",
				conflict.ModulePath,
				conflict.CurrentVersion,
				conflict.ConflictVersion,
				conflict.Reason)
		}
		fmt.Println()
	} else {
		green.Printf("✅ No version conflicts detected\n\n")
	}

	if len(graph.SecurityIssues) > 0 {
		red.Printf("🚨 Security Issues (%d):\n", len(graph.SecurityIssues))
		for _, issue := range graph.SecurityIssues {
			fmt.Printf("  • %s [%s]: %s\n", issue.ID, issue.Severity, issue.Description)
			if issue.FixedIn != "" {
				fmt.Printf("    Fixed in: %s\n", issue.FixedIn)
			}
		}
		fmt.Println()
	} else {
		green.Printf("✅ No known security issues\n\n")
	}

	blue.Printf("📄 License Summary:\n")
	for license, count := range graph.LicensesSummary {
		fmt.Printf("  • %s: %d packages\n", license, count)
	}
	fmt.Println()

	yellow.Printf("💡 Recommendations:\n")
	if len(graph.Conflicts) > 0 {
		fmt.Printf("  • Review and resolve version conflicts\n")
	}
	if len(graph.SecurityIssues) > 0 {
		fmt.Printf("  • Update packages with security vulnerabilities\n")
	}
	if graph.LicensesSummary["Unknown"] > 0 {
		fmt.Printf("  • Review licenses for %d unknown packages\n", graph.LicensesSummary["Unknown"])
	}
	fmt.Printf("  • Consider running 'go mod tidy' to clean up dependencies\n")
	fmt.Printf("  • Use 'goviz doctor' for detailed package health analysis\n")

	return nil
}

func init() {
	analyzeCmd.Flags().StringVarP(&analyzeFormat, "format", "f", "text", "Output format (json, yaml, text, console)")
	analyzeCmd.Flags().StringVarP(&analyzeOutput, "output", "o", "", "Output file (stdout if not specified)")
	analyzeCmd.Flags().BoolVar(&showConflicts, "conflicts", false, "Show only version conflicts")
	analyzeCmd.Flags().BoolVar(&showOutdated, "outdated", false, "Show only outdated packages")
}
