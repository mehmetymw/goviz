package graph

import (
	"golang.org/x/mod/modfile"
)

type Node struct {
	Name     string
	Version  string
	Direct   bool
	Children []*Node
}

type DependencyGraph struct {
	Root            *Node
	AllNodes        map[string]*Node
	ModuleName      string
	ModuleGoVersion string
}

func BuildDependencyGraph(modFile *modfile.File) *DependencyGraph {
	graph := &DependencyGraph{
		AllNodes:   make(map[string]*Node),
		ModuleName: modFile.Module.Mod.Path,
	}

	if modFile.Go != nil {
		graph.ModuleGoVersion = modFile.Go.Version
	}

	root := &Node{
		Name:     modFile.Module.Mod.Path,
		Version:  "main",
		Direct:   true,
		Children: make([]*Node, 0),
	}
	graph.Root = root
	graph.AllNodes[root.Name] = root

	for _, require := range modFile.Require {
		node := &Node{
			Name:     require.Mod.Path,
			Version:  require.Mod.Version,
			Direct:   !require.Indirect,
			Children: make([]*Node, 0),
		}

		graph.AllNodes[node.Name] = node

		if !require.Indirect {
			root.Children = append(root.Children, node)
		}
	}

	return graph
}

func (g *DependencyGraph) GetDirectDependencies() []*Node {
	return g.Root.Children
}

func (g *DependencyGraph) GetAllDependencies() []*Node {
	var deps []*Node
	for name, node := range g.AllNodes {
		if name != g.Root.Name {
			deps = append(deps, node)
		}
	}
	return deps
}

func (g *DependencyGraph) GetDependencyCount() (direct, indirect int) {
	for _, node := range g.AllNodes {
		if node.Name != g.Root.Name {
			if node.Direct {
				direct++
			} else {
				indirect++
			}
		}
	}
	return direct, indirect
}
