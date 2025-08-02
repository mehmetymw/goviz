package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"goviz/pkg/graph"
	"goviz/pkg/parser"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	licensesFormat string
	licensesOutput string
	checkCompat    bool
)

var licensesCmd = &cobra.Command{
	Use:   "licenses [path]",
	Short: "Analyze dependency licenses for compliance",
	Long: `Analyze the licenses of your Go module dependencies.
	
This command:
- Identifies licenses for all dependencies
- Checks for license compatibility issues
- Provides compliance reports
- Flags potentially problematic licenses`,
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

		fmt.Printf("📄 Analyzing dependency licenses...\n")
		modFile, err := parser.ParseGoMod(goModPath)
		if err != nil {
			return fmt.Errorf("failed to parse go.mod: %w", err)
		}

		goSumPath := filepath.Join(absPath, "go.sum")
		enhancedGraph, err := graph.BuildEnhancedDependencyGraph(modFile, goSumPath)
		if err != nil {
			return fmt.Errorf("failed to build enhanced dependency graph: %w", err)
		}

		if err := enhancedGraph.AnalyzeLicenses(); err != nil {
			return fmt.Errorf("failed to analyze licenses: %w", err)
		}

		return generateLicenseReport(enhancedGraph)
	},
}

func generateLicenseReport(graph *graph.EnhancedDependencyGraph) error {
	green := color.New(color.FgGreen, color.Bold)
	yellow := color.New(color.FgYellow, color.Bold)
	red := color.New(color.FgRed, color.Bold)
	blue := color.New(color.FgBlue, color.Bold)

	blue.Printf("📄 License Analysis Report\n")
	blue.Printf("=========================\n\n")

	fmt.Printf("Module: %s\n", graph.ModuleName)
	fmt.Printf("Dependencies analyzed: %d\n\n", len(graph.AllNodes)-1)

	blue.Printf("📊 License Summary:\n")

	type licenseCount struct {
		license string
		count   int
	}
	var licenseCounts []licenseCount
	for license, count := range graph.LicensesSummary {
		licenseCounts = append(licenseCounts, licenseCount{license, count})
	}
	sort.Slice(licenseCounts, func(i, j int) bool {
		return licenseCounts[i].count > licenseCounts[j].count
	})

	for _, lc := range licenseCounts {
		var colorFunc *color.Color
		switch lc.license {
		case "Unknown":
			colorFunc = red
		case "GPL-3.0", "AGPL-3.0", "LGPL-3.0":
			colorFunc = yellow
		default:
			colorFunc = green
		}
		colorFunc.Printf("  • %s: %d packages\n", lc.license, lc.count)
	}
	fmt.Println()

	if checkCompat {
		blue.Printf("🔍 License Compatibility Analysis:\n")

		hasGPL := graph.LicensesSummary["GPL-3.0"] > 0 || graph.LicensesSummary["LGPL-3.0"] > 0 || graph.LicensesSummary["AGPL-3.0"] > 0
		hasUnknown := graph.LicensesSummary["Unknown"] > 0

		if hasGPL {
			yellow.Printf("  ⚠️  Copyleft licenses detected (GPL family)\n")
			fmt.Printf("     • May require source code disclosure\n")
			fmt.Printf("     • Review compatibility with your project license\n")
		}

		if hasUnknown {
			red.Printf("  ❌ Unknown licenses detected\n")
			fmt.Printf("     • Manual license review required\n")
			fmt.Printf("     • Risk of license compliance issues\n")
		}

		if !hasGPL && !hasUnknown {
			green.Printf("  ✅ No obvious license compatibility issues\n")
		}

		fmt.Println()
	}

	blue.Printf("📋 Detailed License Breakdown:\n")

	licensePackages := make(map[string][]string)
	for name, node := range graph.EnhancedNodes {
		if name == graph.Root.Name {
			continue
		}
		license := node.License
		if license == "" {
			license = "Unknown"
		}
		licensePackages[license] = append(licensePackages[license], name)
	}

	var licenses []string
	for license := range licensePackages {
		licenses = append(licenses, license)
	}
	sort.Strings(licenses)

	for _, license := range licenses {
		packages := licensePackages[license]
		sort.Strings(packages)

		var colorFunc *color.Color
		switch license {
		case "Unknown":
			colorFunc = red
		case "GPL-3.0", "AGPL-3.0", "LGPL-3.0":
			colorFunc = yellow
		default:
			colorFunc = green
		}

		colorFunc.Printf("\n%s (%d packages):\n", license, len(packages))
		for _, pkg := range packages {
			fmt.Printf("  • %s\n", pkg)
		}
	}

	fmt.Println()
	yellow.Printf("💡 Recommendations:\n")

	if graph.LicensesSummary["Unknown"] > 0 {
		fmt.Printf("  • Manually review packages with unknown licenses\n")
		fmt.Printf("  • Consider alternative packages with clear licenses\n")
	}

	if graph.LicensesSummary["GPL-3.0"] > 0 || graph.LicensesSummary["LGPL-3.0"] > 0 {
		fmt.Printf("  • Review GPL/LGPL license implications for your project\n")
		fmt.Printf("  • Consult legal team if distributing commercial software\n")
	}

	fmt.Printf("  • Document license compliance in your project\n")
	fmt.Printf("  • Set up automated license scanning in CI/CD\n")
	fmt.Printf("  • Consider using license compatibility tools\n")

	return nil
}

func init() {
	licensesCmd.Flags().StringVarP(&licensesFormat, "format", "f", "text", "Output format (text, json, yaml)")
	licensesCmd.Flags().StringVarP(&licensesOutput, "output", "o", "", "Output file")
	licensesCmd.Flags().BoolVar(&checkCompat, "check-compatibility", true, "Check license compatibility")
}
