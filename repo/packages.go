package repo

import (
	"os"
	"path/filepath"
	"strings"

	parser "github.com/autamus/binoc/repo"
)

// GetChangedPackages returns a slice of parser.Results from the files changed.
func GetChangedPackages(prefixPath, packagesPath string, filepaths []string) (packages []parser.Result, err error) {
	for _, path := range filepaths {
		if strings.Contains(path, packagesPath) {
			result, err := parser.Parse(filepath.Join(prefixPath, path))
			if err != nil && err.Error() != "not a valid package format" {
				return packages, err
			}
			if err == nil {
				packages = append(packages, result)
			}
		}
	}
	return packages, nil
}

// FindAndParse attempts to find a package in the repository and parse it.
func FindAndParse(path, name string) (output parser.Result, err error) {
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		dir := strings.ToLower(filepath.Base(filepath.Dir(path)))
		if dir == strings.ToLower(name) {
			result, err := parser.Parse(path)
			if err != nil && err.Error() != "not a valid package format" {
				return err
			}
			if err == nil {
				output = result
			}
		}
		return nil
	})
	return output, err
}
