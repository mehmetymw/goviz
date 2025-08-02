package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type GoSumEntry struct {
	ModulePath string
	Version    string
	Hash       string
}

func ParseGoSum(path string) (map[string]GoSumEntry, error) {
	entries := make(map[string]GoSumEntry)

	file, err := os.Open(path)
	if err != nil {

		return entries, nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}

		modulePath := parts[0]
		version := parts[1]
		hash := parts[2]

		if !strings.HasSuffix(version, "/go.mod") {
			entries[modulePath+"@"+version] = GoSumEntry{
				ModulePath: modulePath,
				Version:    version,
				Hash:       hash,
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading go.sum: %w", err)
	}

	return entries, nil
}

func GetTransitiveDependencies(goSumEntries map[string]GoSumEntry, directDeps []string) []GoSumEntry {
	directDepMap := make(map[string]bool)
	for _, dep := range directDeps {
		directDepMap[dep] = true
	}

	var transitive []GoSumEntry
	seenModules := make(map[string]bool)

	for _, entry := range goSumEntries {
		if !directDepMap[entry.ModulePath] && !seenModules[entry.ModulePath] {
			transitive = append(transitive, entry)
			seenModules[entry.ModulePath] = true
		}
	}

	return transitive
}
