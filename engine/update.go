package engine

import (
	"errors"
	"strings"

	parser "github.com/autamus/binoc/repo"
)

func findUpdate(packagePath, commitMsg string, packages []parser.Result) (result parser.Result, err error) {
	reverseDeps := make(map[string][]parser.Result)
	for _, pkg := range packages {
		deps := pkg.Package.GetDependencies()
		for _, dep := range deps {
			reverseDeps[dep] = append(reverseDeps[dep], pkg)
		}
	}

	if len(packages) > 0 {
		path := []parser.Result{}
		current := packages[0]
		for {
			upstreams := reverseDeps[current.Package.GetName()]
			if len(upstreams) == 0 {
				return current, nil
			}
			for _, entry := range upstreams {
				if !contains(path, entry.Package.GetName()) {
					path = append(path, current)
					current = entry
					continue
				}
			}
			break
		}

		for _, entry := range path {
			if strings.Contains(commitMsg, strings.ToLower(entry.Package.GetName())) {
				return entry, nil
			}
		}
	}

	return result, errors.New("could not find head package to update")
}

func contains(from []parser.Result, search string) bool {
	for _, a := range from {
		if a.Package.GetName() == search {
			return true
		}
	}
	return false
}
