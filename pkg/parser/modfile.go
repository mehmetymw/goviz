package parser

import (
	"fmt"
	"os"

	"golang.org/x/mod/modfile"
)

func ParseGoMod(path string) (*modfile.File, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read go.mod file: %w", err)
	}

	modFile, err := modfile.Parse(path, data, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to parse go.mod file: %w", err)
	}

	return modFile, nil
}

func GetDirectDependencies(modFile *modfile.File) []string {
	var deps []string

	for _, require := range modFile.Require {

		if !require.Indirect {
			deps = append(deps, require.Mod.Path)
		}
	}

	return deps
}

func GetAllDependencies(modFile *modfile.File) []string {
	var deps []string

	for _, require := range modFile.Require {
		deps = append(deps, require.Mod.Path)
	}

	return deps
}
