package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"goviz/pkg/graph"
	"goviz/pkg/parser"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	securitySeverity string
	securityFormat   string
	securityOutput   string
)

var securityCmd = &cobra.Command{
	Use:   "security [path]",
	Short: "Check dependencies for security vulnerabilities",
	Long: `Check your Go module dependencies for known security vulnerabilities.
	
This command:
- Scans all dependencies for known CVEs
- Reports vulnerability severity levels
- Suggests fixes and updates
- Provides actionable security recommendations`,
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

		fmt.Printf("🔒 Scanning dependencies for security vulnerabilities...\n")
		modFile, err := parser.ParseGoMod(goModPath)
		if err != nil {
			return fmt.Errorf("failed to parse go.mod: %w", err)
		}

		goSumPath := filepath.Join(absPath, "go.sum")
		enhancedGraph, err := graph.BuildEnhancedDependencyGraph(modFile, goSumPath)
		if err != nil {
			return fmt.Errorf("failed to build enhanced dependency graph: %w", err)
		}

		if err := enhancedGraph.CheckSecurity(); err != nil {
			return fmt.Errorf("failed to check security: %w", err)
		}

		return generateSecurityReport(enhancedGraph)
	},
}

func generateSecurityReport(depGraph *graph.EnhancedDependencyGraph) error {
	red := color.New(color.FgRed, color.Bold)
	green := color.New(color.FgGreen, color.Bold)
	yellow := color.New(color.FgYellow, color.Bold)
	blue := color.New(color.FgBlue, color.Bold)

	blue.Printf("🔒 Security Vulnerability Report\n")
	blue.Printf("================================\n\n")

	fmt.Printf("Module: %s\n", depGraph.ModuleName)
	fmt.Printf("Scanned: %d dependencies\n\n", len(depGraph.AllNodes)-1)

	if len(depGraph.SecurityIssues) == 0 {
		green.Printf("✅ No known security vulnerabilities found!\n\n")

		blue.Printf("🛡️  Security Recommendations:\n")
		fmt.Printf("  • Keep dependencies up to date\n")
		fmt.Printf("  • Regularly run security scans\n")
		fmt.Printf("  • Monitor security advisories\n")
		fmt.Printf("  • Use 'go mod tidy' to remove unused dependencies\n")

		return nil
	}

	severityCount := make(map[string]int)
	severityIssues := make(map[string][]any)

	for _, issue := range depGraph.SecurityIssues {
		severityCount[issue.Severity]++
		severityIssues[issue.Severity] = append(severityIssues[issue.Severity], any(issue))
	}

	red.Printf("🚨 Found %d security issues:\n", len(depGraph.SecurityIssues))
	for severity, count := range severityCount {
		var colorFunc func(a ...any) string
		switch severity {
		case "CRITICAL":
			colorFunc = red.Sprint
		case "HIGH":
			colorFunc = red.Sprint
		case "MEDIUM":
			colorFunc = yellow.Sprint
		case "LOW":
			colorFunc = green.Sprint
		default:
			colorFunc = fmt.Sprint
		}
		fmt.Printf("  • %s: %d\n", colorFunc("%s", severity), count)
	}
	fmt.Println()

	severityOrder := []string{"CRITICAL", "HIGH", "MEDIUM", "LOW"}
	for _, severity := range severityOrder {
		issues := severityIssues[severity]
		if len(issues) == 0 {
			continue
		}

		var colorFunc *color.Color
		switch severity {
		case "CRITICAL", "HIGH":
			colorFunc = red
		case "MEDIUM":
			colorFunc = yellow
		default:
			colorFunc = green
		}
		colorFunc.Printf("%s Severity Issues:\n", severity)
		for i, issueInterface := range issues {
			issue := issueInterface.(graph.SecurityIssue)
			fmt.Printf("  %d. %s\n", i+1, issue.ID)
			fmt.Printf("     Description: %s\n", issue.Description)
			if issue.FixedIn != "" {
				fmt.Printf("     Fixed in: %s\n", issue.FixedIn)
			} else {
				fmt.Printf("     Fixed in: N/A\n")
			}
			fmt.Println()
		}
	}

	yellow.Printf("🔧 Recommended Actions:\n")
	if severityCount["CRITICAL"] > 0 || severityCount["HIGH"] > 0 {
		fmt.Printf("  🚨 URGENT: Update packages with CRITICAL/HIGH severity issues immediately\n")
	}
	if severityCount["MEDIUM"] > 0 {
		fmt.Printf("  ⚠️  Plan updates for MEDIUM severity issues in next release\n")
	}
	if severityCount["LOW"] > 0 {
		fmt.Printf("  📝 Consider updating LOW severity issues when convenient\n")
	}

	fmt.Printf("  • Run 'go get -u' to update dependencies\n")
	fmt.Printf("  • Review and test updates in development environment\n")
	fmt.Printf("  • Set up automated security scanning in CI/CD\n")

	if severityCount["CRITICAL"] > 0 || severityCount["HIGH"] > 0 {
		fmt.Printf("\n")
		red.Printf("❌ Security scan failed due to high-severity vulnerabilities\n")
		os.Exit(1)
	}

	return nil
}

func init() {
	securityCmd.Flags().StringVarP(&securitySeverity, "severity", "s", "", "Filter by severity (CRITICAL, HIGH, MEDIUM, LOW)")
	securityCmd.Flags().StringVarP(&securityFormat, "format", "f", "text", "Output format (text, json, yaml)")
	securityCmd.Flags().StringVarP(&securityOutput, "output", "o", "", "Output file")
}
