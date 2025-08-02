package graph

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"goviz/pkg/parser"

	"golang.org/x/mod/modfile"
)

type EnhancedNode struct {
	*Node
	Hash            string
	Transitive      []*EnhancedNode
	Conflicts       []VersionConflict
	SecurityIssues  []SecurityIssue
	License         string
	LastUpdate      time.Time
	IsOutdated      bool
	UpdateAvailable string
}

type VersionConflict struct {
	ModulePath      string
	CurrentVersion  string
	ConflictVersion string
	Reason          string
}

type SecurityIssue struct {
	ID          string
	Severity    string
	Description string
	FixedIn     string
}

type EnhancedDependencyGraph struct {
	*DependencyGraph
	EnhancedNodes   map[string]*EnhancedNode
	GoSumEntries    map[string]parser.GoSumEntry
	Conflicts       []VersionConflict
	SecurityIssues  []SecurityIssue
	TotalSize       int64
	BuildTime       time.Duration
	LicensesSummary map[string]int
}

func BuildEnhancedDependencyGraph(modFile *modfile.File, goSumPath string) (*EnhancedDependencyGraph, error) {

	basicGraph := BuildDependencyGraph(modFile)

	goSumEntries, err := parser.ParseGoSum(goSumPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse go.sum: %w", err)
	}

	enhancedGraph := &EnhancedDependencyGraph{
		DependencyGraph: basicGraph,
		EnhancedNodes:   make(map[string]*EnhancedNode),
		GoSumEntries:    goSumEntries,
		LicensesSummary: make(map[string]int),
	}

	for name, node := range basicGraph.AllNodes {
		enhancedNode := &EnhancedNode{
			Node:           node,
			Transitive:     make([]*EnhancedNode, 0),
			Conflicts:      make([]VersionConflict, 0),
			SecurityIssues: make([]SecurityIssue, 0),
		}

		key := name + "@" + node.Version
		if entry, exists := goSumEntries[key]; exists {
			enhancedNode.Hash = entry.Hash
		}

		enhancedGraph.EnhancedNodes[name] = enhancedNode
	}

	directDeps := parser.GetDirectDependencies(modFile)
	transitiveDeps := parser.GetTransitiveDependencies(goSumEntries, directDeps)

	for _, transDep := range transitiveDeps {
		if _, exists := enhancedGraph.EnhancedNodes[transDep.ModulePath]; !exists {
			node := &Node{
				Name:     transDep.ModulePath,
				Version:  transDep.Version,
				Direct:   false,
				Children: make([]*Node, 0),
			}

			enhancedNode := &EnhancedNode{
				Node:           node,
				Hash:           transDep.Hash,
				Transitive:     make([]*EnhancedNode, 0),
				Conflicts:      make([]VersionConflict, 0),
				SecurityIssues: make([]SecurityIssue, 0),
			}

			enhancedGraph.EnhancedNodes[transDep.ModulePath] = enhancedNode
			enhancedGraph.AllNodes[transDep.ModulePath] = node
		}
	}

	return enhancedGraph, nil
}

func (g *EnhancedDependencyGraph) DetectVersionConflicts() {
	versionMap := make(map[string][]string)

	for _, entry := range g.GoSumEntries {
		versionMap[entry.ModulePath] = append(versionMap[entry.ModulePath], entry.Version)
	}

	for modulePath, versions := range versionMap {
		if len(versions) > 1 {
			sort.Strings(versions)
			for i := 0; i < len(versions)-1; i++ {
				conflict := VersionConflict{
					ModulePath:      modulePath,
					CurrentVersion:  versions[len(versions)-1],
					ConflictVersion: versions[i],
					Reason:          "Multiple versions in go.sum",
				}
				g.Conflicts = append(g.Conflicts, conflict)

				if node, exists := g.EnhancedNodes[modulePath]; exists {
					node.Conflicts = append(node.Conflicts, conflict)
				}
			}
		}
	}
}

func (g *EnhancedDependencyGraph) AnalyzeLicenses() error {

	knownLicenses := map[string]string{
		"github.com/spf13/cobra":               "Apache-2.0",
		"github.com/spf13/pflag":               "BSD-3-Clause",
		"github.com/awalterschulze/gographviz": "Apache-2.0",
		"github.com/inconshreveable/mousetrap": "Apache-2.0",
		"golang.org/x/mod":                     "BSD-3-Clause",
		"gopkg.in/yaml.v3":                     "Apache-2.0",
		"github.com/google/licensecheck":       "BSD-3-Clause",
		"github.com/fatih/color":               "MIT",
	}

	for name, node := range g.EnhancedNodes {
		if license, exists := knownLicenses[name]; exists {
			node.License = license
			g.LicensesSummary[license]++
		} else {

			if strings.Contains(name, "golang.org/x/") {
				node.License = "BSD-3-Clause"
				g.LicensesSummary["BSD-3-Clause"]++
			} else if strings.Contains(name, "github.com/mattn/") {
				node.License = "MIT"
				g.LicensesSummary["MIT"]++
			} else {
				node.License = "Unknown"
				g.LicensesSummary["Unknown"]++
			}
		}
	}

	return nil
}

func (g *EnhancedDependencyGraph) CheckSecurity() error {

	vulnerablePatterns := map[string]SecurityIssue{

		"github.com/gin-gonic/gin": {
			ID:          "GHSA-example",
			Severity:    "MEDIUM",
			Description: "Check for latest version with security fixes",
			FixedIn:     "v1.9.1+",
		},
		"github.com/gorilla/websocket": {
			ID:          "CVE-2023-example",
			Severity:    "HIGH",
			Description: "WebSocket vulnerability in older versions",
			FixedIn:     "v1.5.0+",
		},
	}

	for name, node := range g.EnhancedNodes {
		if name == g.Root.Name {
			continue
		}

		if issue, exists := vulnerablePatterns[name]; exists {

			if strings.Contains(node.Version, "v1.8") ||
				strings.Contains(node.Version, "v1.7") ||
				strings.Contains(node.Version, "v1.4") {
				node.SecurityIssues = append(node.SecurityIssues, issue)
				g.SecurityIssues = append(g.SecurityIssues, issue)
			}
		}

		if strings.Contains(node.Version, "dev") ||
			strings.Contains(node.Version, "alpha") ||
			strings.Contains(node.Version, "beta") ||
			strings.Contains(node.Version, "rc") ||
			strings.Contains(node.Version, "snapshot") {
			issue := SecurityIssue{
				ID:          "DEV-VERSION",
				Severity:    "LOW",
				Description: "Development version detected in dependencies",
				FixedIn:     "Use stable release version",
			}
			node.SecurityIssues = append(node.SecurityIssues, issue)
			g.SecurityIssues = append(g.SecurityIssues, issue)
		}

		if strings.Contains(node.Version, "20161208") ||
			strings.Contains(node.Version, "20170") ||
			strings.Contains(node.Version, "20180") {
			issue := SecurityIssue{
				ID:          "OLD-VERSION",
				Severity:    "MEDIUM",
				Description: "Very old package version may have security vulnerabilities",
				FixedIn:     "Update to latest version",
			}
			node.SecurityIssues = append(node.SecurityIssues, issue)
			g.SecurityIssues = append(g.SecurityIssues, issue)
		}

		insecurePatterns := []string{
			"crypto/md5",
			"crypto/sha1",
			"net/http/httputil",
		}

		for _, pattern := range insecurePatterns {
			if strings.Contains(name, pattern) {
				issue := SecurityIssue{
					ID:          "INSECURE-CRYPTO",
					Severity:    "HIGH",
					Description: "Package uses insecure cryptographic functions",
					FixedIn:     "Use secure alternatives (SHA-256, bcrypt, etc.)",
				}
				node.SecurityIssues = append(node.SecurityIssues, issue)
				g.SecurityIssues = append(g.SecurityIssues, issue)
			}
		}

		if node.Version == "" || node.Version == "v0.0.0" {
			issue := SecurityIssue{
				ID:          "NO-VERSION",
				Severity:    "LOW",
				Description: "Package without proper versioning detected",
				FixedIn:     "Use properly versioned packages",
			}
			node.SecurityIssues = append(node.SecurityIssues, issue)
			g.SecurityIssues = append(g.SecurityIssues, issue)
		}
	}

	return nil
}

func (g *EnhancedDependencyGraph) GetStatistics() map[string]any {
	direct, indirect := g.GetDependencyCount()
	transitive := len(g.GoSumEntries) - direct - indirect

	stats := map[string]any{
		"total_dependencies":      len(g.AllNodes) - 1,
		"direct_dependencies":     direct,
		"indirect_dependencies":   indirect,
		"transitive_dependencies": transitive,
		"version_conflicts":       len(g.Conflicts),
		"security_issues":         len(g.SecurityIssues),
		"unique_licenses":         len(g.LicensesSummary),
		"licenses_breakdown":      g.LicensesSummary,
	}

	return stats
}
