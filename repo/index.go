package repo

import (
	"os"
	"path/filepath"
	"strings"

	binoc "github.com/autamus/binoc/repo"
)

func IndexReverseDependencies(parser binoc.Repo, path, packagesPath string) (result map[string][]string, err error) {
	// Iterate through and parse packages
	output := make(chan binoc.Result, 20)
	result = make(map[string][]string)
	go parser.ParseDir(filepath.Join(path, packagesPath), output)

	// Construct reverse dependency map by mapping dependencies to the list of
	// apps that depend on them.
	for app := range output {
		for _, dependency := range app.Package.GetDependencies() {
			result[dependency] = append(result[dependency], app.Package.GetName())
		}
	}
	return result, nil
}

func IndexPackageContainerDeps(prefixPath, defaultEnvPath, containersPath string) (result map[string][]string, err error) {
	// Initialize empty reverse depencency map
	result = make(map[string][]string)

	// Walk through and parse spack env containers.
	location := filepath.Join(prefixPath, containersPath)
	err = filepath.Walk(location, func(path string, info os.FileInfo, err error) error {
		// Check if file is a Spack Env Yaml Spec
		match, _ := filepath.Match("spack.yaml", filepath.Base(path))
		if match {
			// Find container name
			containerName := filepath.Base(filepath.Dir(path))
			// Parse a spack environment into a struct.
			container, err := ParseSpackEnv(defaultEnvPath, path)
			if err != nil {
				return err
			}

			// Add container as dependency of package.
			for _, spec := range container.Spack.Specs {
				// Record the end of the dependency name versus version/variant info.
				end := strings.IndexFunc(spec, versend)
				if end > 0 {
					spec = spec[:end]
				}
				result[spec] = append(result[spec], containerName)
			}
		}
		return nil
	})

	return result, err
}
