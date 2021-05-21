package repo

import (
	"path/filepath"
	"strings"
)

// GetChangedContainers returns a slice of parser.Results from the files changed.
func GetChangedContainers(prefixPath, containersPath string, filepaths []string) (containers map[string]bool, err error) {
	containers = make(map[string]bool)
	for _, path := range filepaths {
		if strings.Contains(path, containersPath) {
			containerName := filepath.Base(filepath.Dir(path))
			containers[containerName] = true
		}
	}
	return containers, nil
}
