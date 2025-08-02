package output

import (
	"fmt"
	"os"
	"strings"

	"goviz/pkg/graph"

	"github.com/awalterschulze/gographviz"
)

func GenerateDOT(depGraph *graph.DependencyGraph, outputFile string) error {

	graphAst, err := gographviz.ParseString(`digraph G {}`)
	if err != nil {
		return fmt.Errorf("failed to create graph: %w", err)
	}

	graph := gographviz.NewGraph()
	if err = gographviz.Analyse(graphAst, graph); err != nil {
		return fmt.Errorf("failed to analyze graph: %w", err)
	}

	if err := graph.SetName("DependencyGraph"); err != nil {
		return fmt.Errorf("failed to set graph name: %w", err)
	}
	if err := graph.SetDir(true); err != nil {
		return fmt.Errorf("failed to set graph direction: %w", err)
	}

	if err := graph.AddAttr("DependencyGraph", "rankdir", "TB"); err != nil {
		return fmt.Errorf("failed to add rankdir attribute: %w", err)
	}
	if err := graph.AddAttr("DependencyGraph", "node", "[shape=box,style=filled]"); err != nil {
		return fmt.Errorf("failed to add node attribute: %w", err)
	}

	rootNodeName := sanitizeNodeName(depGraph.Root.Name)
	if err := graph.AddNode("DependencyGraph", rootNodeName, map[string]string{
		"label":     fmt.Sprintf("\"%s\\n(main)\"", depGraph.Root.Name),
		"fillcolor": "lightblue",
		"style":     "filled",
	}); err != nil {
		return fmt.Errorf("failed to add root node: %w", err)
	}

	for _, node := range depGraph.GetAllDependencies() {
		nodeName := sanitizeNodeName(node.Name)
		color := "lightgray"
		if node.Direct {
			color = "lightgreen"
		}

		label := fmt.Sprintf("\"%s\\n%s\"", node.Name, node.Version)
		if err := graph.AddNode("DependencyGraph", nodeName, map[string]string{
			"label":     label,
			"fillcolor": color,
			"style":     "filled",
		}); err != nil {
			return fmt.Errorf("failed to add node %s: %w", node.Name, err)
		}

		if node.Direct {
			if err := graph.AddEdge(rootNodeName, nodeName, true, nil); err != nil {
				return fmt.Errorf("failed to add edge from %s to %s: %w", depGraph.Root.Name, node.Name, err)
			}
		}
	}

	dotContent := graph.String()
	if err := os.WriteFile(outputFile, []byte(dotContent), 0644); err != nil {
		return fmt.Errorf("failed to write DOT file: %w", err)
	}

	fmt.Printf("DOT file generated: %s\n", outputFile)
	fmt.Printf("To visualize: dot -Tpng %s -o depgraph.png\n", outputFile)

	return nil
}

func sanitizeNodeName(name string) string {

	sanitized := strings.ReplaceAll(name, "/", "_")
	sanitized = strings.ReplaceAll(sanitized, ".", "_")
	sanitized = strings.ReplaceAll(sanitized, "-", "_")
	return fmt.Sprintf("\"%s\"", sanitized)
}
