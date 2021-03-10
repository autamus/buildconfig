package repo

import (
	"os"
	"path/filepath"
	"strings"

	parser "github.com/autamus/binoc/repo"
)

// GetPackages returns a slice of parser.Results from the files changed.
func GetPackages(prefixPath, packagesPath string, filepaths []string) (packages []parser.Result, err error) {
	for _, path := range filepaths {
		if strings.Contains(path, packagesPath) {
			result, err := parser.Parse(filepath.Join(prefixPath, path))
			if err != nil {
				return packages, err
			}
			packages = append(packages, result)
		}
	}
	return packages, nil
}

// FindAndParse attempts to find a package in the repository and parse it.
func FindAndParse(path, name string) (output parser.Result, err error) {
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		dir := strings.ToLower(filepath.Base(filepath.Dir(path)))
		if dir == strings.ToLower(name) {
			output, err = parser.Parse(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return output, err
}
