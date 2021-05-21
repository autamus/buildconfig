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
		result[ToHyphenCase(currentPack)] = true

		// Deqeue current after adding to map.
		packages = packages[1:]
		// Append all reverse dependencies to packages queue.
		packages = append(packages, packageDeps[currentPack]...)
	}
	return result
}
