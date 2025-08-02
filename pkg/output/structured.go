package output

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"goviz/pkg/graph"

	"gopkg.in/yaml.v3"
)

type DependencyReport struct {
	Metadata        ReportMetadata          `json:"metadata" yaml:"metadata"`
	Module          ModuleInfo              `json:"module" yaml:"module"`
	Statistics      map[string]any          `json:"statistics" yaml:"statistics"`
	Dependencies    []DependencyInfo        `json:"dependencies" yaml:"dependencies"`
	Conflicts       []graph.VersionConflict `json:"conflicts,omitempty" yaml:"conflicts,omitempty"`
	SecurityIssues  []graph.SecurityIssue   `json:"security_issues,omitempty" yaml:"security_issues,omitempty"`
	LicensesSummary map[string]int          `json:"licenses_summary" yaml:"licenses_summary"`
}

type ReportMetadata struct {
	GeneratedAt time.Time `json:"generated_at" yaml:"generated_at"`
	Tool        string    `json:"tool" yaml:"tool"`
	Version     string    `json:"version" yaml:"version"`
}

type ModuleInfo struct {
	Name      string `json:"name" yaml:"name"`
	GoVersion string `json:"go_version,omitempty" yaml:"go_version,omitempty"`
	Path      string `json:"path" yaml:"path"`
}

type DependencyInfo struct {
	Name            string                  `json:"name" yaml:"name"`
	Version         string                  `json:"version" yaml:"version"`
	Direct          bool                    `json:"direct" yaml:"direct"`
	Hash            string                  `json:"hash,omitempty" yaml:"hash,omitempty"`
	License         string                  `json:"license,omitempty" yaml:"license,omitempty"`
	Conflicts       []graph.VersionConflict `json:"conflicts,omitempty" yaml:"conflicts,omitempty"`
	SecurityIssues  []graph.SecurityIssue   `json:"security_issues,omitempty" yaml:"security_issues,omitempty"`
	IsOutdated      bool                    `json:"is_outdated,omitempty" yaml:"is_outdated,omitempty"`
	UpdateAvailable string                  `json:"update_available,omitempty" yaml:"update_available,omitempty"`
}

func GenerateJSON(depGraph *graph.EnhancedDependencyGraph, outputFile, projectPath string) error {
	report := buildDependencyReport(depGraph, projectPath)

	jsonData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if outputFile == "" {
		fmt.Print(string(jsonData))
		return nil
	}

	if err := os.WriteFile(outputFile, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	fmt.Printf("JSON report generated: %s\n", outputFile)
	return nil
}

func GenerateYAML(depGraph *graph.EnhancedDependencyGraph, outputFile, projectPath string) error {
	report := buildDependencyReport(depGraph, projectPath)

	yamlData, err := yaml.Marshal(report)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	if outputFile == "" {
		fmt.Print(string(yamlData))
		return nil
	}

	if err := os.WriteFile(outputFile, yamlData, 0644); err != nil {
		return fmt.Errorf("failed to write YAML file: %w", err)
	}

	fmt.Printf("YAML report generated: %s\n", outputFile)
	return nil
}

func buildDependencyReport(depGraph *graph.EnhancedDependencyGraph, projectPath string) DependencyReport {
	var dependencies []DependencyInfo

	for name, enhancedNode := range depGraph.EnhancedNodes {
		if name == depGraph.Root.Name {
			continue
		}

		dep := DependencyInfo{
			Name:            enhancedNode.Name,
			Version:         enhancedNode.Version,
			Direct:          enhancedNode.Direct,
			Hash:            enhancedNode.Hash,
			License:         enhancedNode.License,
			Conflicts:       enhancedNode.Conflicts,
			SecurityIssues:  enhancedNode.SecurityIssues,
			IsOutdated:      enhancedNode.IsOutdated,
			UpdateAvailable: enhancedNode.UpdateAvailable,
		}
		dependencies = append(dependencies, dep)
	}

	return DependencyReport{
		Metadata: ReportMetadata{
			GeneratedAt: time.Now(),
			Tool:        "goviz",
			Version:     "v0.1.0",
		},
		Module: ModuleInfo{
			Name:      depGraph.ModuleName,
			GoVersion: depGraph.ModuleGoVersion,
			Path:      projectPath,
		},
		Statistics:      depGraph.GetStatistics(),
		Dependencies:    dependencies,
		Conflicts:       depGraph.Conflicts,
		SecurityIssues:  depGraph.SecurityIssues,
		LicensesSummary: depGraph.LicensesSummary,
	}
}
