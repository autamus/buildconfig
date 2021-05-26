package engine

import (
	parser "github.com/autamus/binoc/repo"
)

func GetAllPackageBuilds(changed []parser.Result, packageDeps map[string][]string) (result map[string]bool) {
	result = make(map[string]bool)
	packages := []string{}

	for _, pack := range changed {
		packages = append(packages, pack.Package.GetName())
	}

	for len(packages) > 0 {
		currentPack := packages[0]
		// Deqeue current after adding to map.
		packages = packages[1:]
		// Prevent loop in the case of circular dependencies.
		if result[ToHyphenCase(currentPack)] {
			continue
		}
		result[ToHyphenCase(currentPack)] = true
		// Append all reverse dependencies to packages queue.
		packages = append(packages, packageDeps[ToHyphenCase(currentPack)]...)
	}
	return result
}
