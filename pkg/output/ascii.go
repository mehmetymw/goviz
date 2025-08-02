package output

import (
	"fmt"
	"sort"
	"strings"

	"goviz/pkg/graph"
)

func GenerateASCIITree(depGraph *graph.DependencyGraph) error {
	fmt.Printf("Dependency Graph for: %s\n", depGraph.ModuleName)

	if depGraph.ModuleGoVersion != "" {
		fmt.Printf("Go Version: %s\n", depGraph.ModuleGoVersion)
	}

	direct, indirect := depGraph.GetDependencyCount()
	fmt.Printf("Dependencies: %d direct, %d indirect\n", direct, indirect)
	fmt.Println()

	fmt.Printf("%s (main)\n", depGraph.Root.Name)

	directDeps := depGraph.GetDirectDependencies()
	if len(directDeps) == 0 {
		fmt.Println("└── (no dependencies)")
		return nil
	}

	sort.Slice(directDeps, func(i, j int) bool {
		return directDeps[i].Name < directDeps[j].Name
	})

	for i, dep := range directDeps {
		isLast := i == len(directDeps)-1
		printNode(dep, "", isLast)
	}

	allDeps := depGraph.GetAllDependencies()
	var indirectDeps []*graph.Node
	for _, dep := range allDeps {
		if !dep.Direct {
			indirectDeps = append(indirectDeps, dep)
		}
	}

	if len(indirectDeps) > 0 {
		fmt.Printf("\nIndirect dependencies (%d):\n", len(indirectDeps))

		sort.Slice(indirectDeps, func(i, j int) bool {
			return indirectDeps[i].Name < indirectDeps[j].Name
		})

		for i, dep := range indirectDeps {
			isLast := i == len(indirectDeps)-1
			var prefix string
			if isLast {
				prefix = "└── "
			} else {
				prefix = "├── "
			}
			fmt.Printf("%s%s (%s)\n", prefix, dep.Name, dep.Version)
		}
	}

	return nil
}

func printNode(node *graph.Node, prefix string, isLast bool) {
	var connector, childPrefix string

	if isLast {
		connector = "└── "
		childPrefix = prefix + "    "
	} else {
		connector = "├── "
		childPrefix = prefix + "│   "
	}

	fmt.Printf("%s%s%s (%s)\n", prefix, connector, node.Name, node.Version)

	for i, child := range node.Children {
		isChildLast := i == len(node.Children)-1
		printNode(child, childPrefix, isChildLast)
	}
}

func GenerateASCIITreeCompact(depGraph *graph.DependencyGraph) error {
	fmt.Printf("Module: %s\n", depGraph.ModuleName)

	direct, indirect := depGraph.GetDependencyCount()
	fmt.Printf("Dependencies: %d direct, %d indirect\n\n", direct, indirect)

	directDeps := depGraph.GetDirectDependencies()
	allDeps := depGraph.GetAllDependencies()

	var indirectDeps []*graph.Node
	for _, dep := range allDeps {
		if !dep.Direct {
			indirectDeps = append(indirectDeps, dep)
		}
	}

	if len(directDeps) > 0 {
		fmt.Println("Direct dependencies:")
		sort.Slice(directDeps, func(i, j int) bool {
			return directDeps[i].Name < directDeps[j].Name
		})

		for _, dep := range directDeps {
			fmt.Printf("  • %s (%s)\n", dep.Name, dep.Version)
		}
		fmt.Println()
	}

	if len(indirectDeps) > 0 {
		fmt.Println("Indirect dependencies:")
		sort.Slice(indirectDeps, func(i, j int) bool {
			return indirectDeps[i].Name < indirectDeps[j].Name
		})

		grouped := make(map[string][]*graph.Node)
		for _, dep := range indirectDeps {
			parts := strings.Split(dep.Name, "/")
			key := parts[0]
			if len(parts) > 1 {
				key = strings.Join(parts[:2], "/")
			}
			grouped[key] = append(grouped[key], dep)
		}

		var keys []string
		for k := range grouped {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, key := range keys {
			deps := grouped[key]
			if len(deps) == 1 {
				fmt.Printf("  • %s (%s)\n", deps[0].Name, deps[0].Version)
			} else {
				fmt.Printf("  • %s/... (%d packages)\n", key, len(deps))
				for _, dep := range deps {
					fmt.Printf("    - %s (%s)\n", dep.Name, dep.Version)
				}
			}
		}
	}

	return nil
}
