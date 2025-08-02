package output

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"goviz/pkg/graph"
)

func GeneratePNG(depGraph *graph.EnhancedDependencyGraph, outputFile string) error {

	if err := checkGraphvizInstalled(); err != nil {
		return err
	}

	tempDotFile := "temp_depgraph.dot"
	defer os.Remove(tempDotFile)

	if err := GenerateEnhancedDOT(depGraph, tempDotFile); err != nil {
		return fmt.Errorf("failed to generate DOT file: %w", err)
	}

	if outputFile == "" {
		outputFile = "depgraph.png"
	}

	cmd := exec.Command("dot", "-Tpng", tempDotFile, "-o", outputFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to generate PNG: %w\nOutput: %s", err, string(output))
	}

	fmt.Printf("PNG diagram generated: %s\n", outputFile)
	fmt.Printf("Generated using Graphviz dot command\n")

	return nil
}

func GenerateEnhancedDOT(depGraph *graph.EnhancedDependencyGraph, outputFile string) error {

	if err := GenerateDOT(depGraph.DependencyGraph, outputFile); err != nil {
		return err
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		return fmt.Errorf("failed to read DOT file: %w", err)
	}

	enhancedContent := enhanceDOTContent(string(content), depGraph)

	if err := os.WriteFile(outputFile, []byte(enhancedContent), 0644); err != nil {
		return fmt.Errorf("failed to write enhanced DOT file: %w", err)
	}

	return nil
}

func enhanceDOTContent(content string, depGraph *graph.EnhancedDependencyGraph) string {
	lines := strings.Split(content, "\n")
	var enhancedLines []string

	for _, line := range lines {
		if strings.Contains(line, "digraph DependencyGraph") {
			enhancedLines = append(enhancedLines, line)

			enhancedLines = append(enhancedLines, "    graph [fontname=\"Arial\", fontsize=12];")
			enhancedLines = append(enhancedLines, "    node [fontname=\"Arial\", fontsize=10, shape=box];")
			enhancedLines = append(enhancedLines, "    edge [fontname=\"Arial\", fontsize=8];")
		} else if strings.Contains(line, "rankdir=TB") {
			enhancedLines = append(enhancedLines, line)

			enhancedLines = append(enhancedLines, "    subgraph cluster_legend {")
			enhancedLines = append(enhancedLines, "        label=\"Legend\";")
			enhancedLines = append(enhancedLines, "        style=filled;")
			enhancedLines = append(enhancedLines, "        color=lightgrey;")
			enhancedLines = append(enhancedLines, "        legend_main [label=\"Main Module\", fillcolor=lightblue, style=filled];")
			enhancedLines = append(enhancedLines, "        legend_direct [label=\"Direct Dependency\", fillcolor=lightgreen, style=filled];")
			enhancedLines = append(enhancedLines, "        legend_indirect [label=\"Indirect Dependency\", fillcolor=lightgray, style=filled];")
			enhancedLines = append(enhancedLines, "        legend_security [label=\"Security Issue\", fillcolor=red, style=filled];")
			enhancedLines = append(enhancedLines, "    }")
		} else {

			if strings.Contains(line, "[ fillcolor=") && strings.Contains(line, "label=") {
				enhancedLine := enhanceNodeDefinition(line, depGraph)
				enhancedLines = append(enhancedLines, enhancedLine)
			} else {
				enhancedLines = append(enhancedLines, line)
			}
		}
	}

	return strings.Join(enhancedLines, "\n")
}

func enhanceNodeDefinition(line string, depGraph *graph.EnhancedDependencyGraph) string {

	parts := strings.Fields(line)
	if len(parts) == 0 {
		return line
	}

	nodeName := strings.Trim(parts[0], "\"")

	var enhancedNode *graph.EnhancedNode
	for name, node := range depGraph.EnhancedNodes {
		sanitizedName := sanitizeNodeName(name)
		if strings.Contains(sanitizedName, nodeName) || strings.Contains(nodeName, strings.ReplaceAll(name, "/", "_")) {
			enhancedNode = node
			break
		}
	}

	if enhancedNode == nil {
		return line
	}

	if len(enhancedNode.SecurityIssues) > 0 {
		line = strings.ReplaceAll(line, "fillcolor=lightgreen", "fillcolor=red")
		line = strings.ReplaceAll(line, "fillcolor=lightgray", "fillcolor=orange")
	}

	if enhancedNode.License != "" {

		labelStart := strings.Index(line, "label=\"")
		labelEnd := strings.Index(line[labelStart+7:], "\"")
		if labelStart >= 0 && labelEnd >= 0 {
			oldLabel := line[labelStart+7 : labelStart+7+labelEnd]
			newLabel := oldLabel + "\\n" + enhancedNode.License
			if len(enhancedNode.SecurityIssues) > 0 {
				newLabel += "\\n⚠ Security Issues"
			}
			if len(enhancedNode.Conflicts) > 0 {
				newLabel += "\\n⚡ Version Conflicts"
			}
			line = strings.ReplaceAll(line, oldLabel, newLabel)
		}
	}

	return line
}

func checkGraphvizInstalled() error {
	cmd := exec.Command("dot", "-V")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Graphviz not found. Please install Graphviz:\n" +
			"  Ubuntu/Debian: sudo apt-get install graphviz\n" +
			"  macOS: brew install graphviz\n" +
			"  Windows: Download from https://graphviz.org/download/")
	}
	return nil
}

func GenerateSVG(depGraph *graph.EnhancedDependencyGraph, outputFile string) error {
	if err := checkGraphvizInstalled(); err != nil {
		return err
	}

	tempDotFile := "temp_depgraph.dot"
	defer os.Remove(tempDotFile)

	if err := GenerateEnhancedDOT(depGraph, tempDotFile); err != nil {
		return fmt.Errorf("failed to generate DOT file: %w", err)
	}

	if outputFile == "" {
		outputFile = "depgraph.svg"
	}

	cmd := exec.Command("dot", "-Tsvg", tempDotFile, "-o", outputFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to generate SVG: %w\nOutput: %s", err, string(output))
	}

	fmt.Printf("SVG diagram generated: %s\n", outputFile)
	return nil
}
