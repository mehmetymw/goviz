package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	Version   = "dev"
	Commit    = "unknown"
	BuildTime = "unknown"
)

var rootCmd = &cobra.Command{
	Use:     "goviz",
	Version: Version,
	Short:   "A comprehensive Go dependency analysis and visualization tool",
	Long: `goviz is a production-ready CLI tool that analyzes Go module dependencies.

Features:
• Dependency visualization (DOT, PNG, SVG, ASCII tree)
• Comprehensive analysis (conflicts, licenses, health)
• Multiple output formats (JSON, YAML for CI/CD)
• License compliance checking
• Dependency health assessment
• Security framework integration`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(analyzeCmd)
	rootCmd.AddCommand(licensesCmd)
	rootCmd.AddCommand(doctorCmd)
	rootCmd.AddCommand(securityCmd)
}

func SetVersionInfo(version, commit, buildTime string) {
	Version = version
	Commit = commit
	BuildTime = buildTime
	rootCmd.Version = fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, buildTime)
}
