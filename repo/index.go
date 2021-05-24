package repo

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	binoc "github.com/autamus/binoc/repo"
	"gopkg.in/yaml.v2"
)

func IndexReverseDependencies(path, packagesPath string) (result map[string][]string, err error) {
	// Iterate through and parse packages
	output := make(chan binoc.Result, 20)
	result = make(map[string][]string)
	go binoc.ParseDir(filepath.Join(path, packagesPath), output)

	// Construct reverse dependency map by mapping dependencies to the list of
	// apps that depend on them.
	for app := range output {
		for _, dependency := range app.Package.GetDependencies() {
			result[dependency] = append(result[dependency], app.Package.GetName())
		}
	}
	return result, nil
}

type SpackEnv struct {
	Specs []string `yaml:"specs"`
}

func IndexPackageContainerDeps(prefixPath, containersPath string) (result map[string][]string, err error) {
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
			// Create a spack struct.
			container := struct {
				Spack SpackEnv `yaml:"spack"`
			}{}

			// Read file contents into a string.
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			// Unmarshal spack.yaml environment file into a struct.
			err = yaml.Unmarshal([]byte(content), &container)
			if err != nil {
				return err
			}

			// Add container as dependency of package.
			for _, spec := range container.Spack.Specs {
				// Record the end of the depedency name versus version/variant info.
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

// versend returns true at the end of the name of a dependency
func versend(input rune) bool {
	for _, c := range []rune{'@', '~', '+'} {
		if input == c {
			return true
		}
	}
	return false
}
