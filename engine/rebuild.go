package engine

import (
	"errors"
	"path/filepath"
	"strings"

	parser "github.com/autamus/binoc/repo"
	"github.com/autamus/buildconfig/repo"
)

func verifyRebuild(path, packagesPath, commitMsg string, packages []parser.Result) (result parser.Result, err error) {
	words := strings.Fields(commitMsg)
	i := 0
	for _, word := range words {
		i++
		if strings.Contains(word, "rebuild") {
			break
		}
	}
	if len(words) <= i+1 {
		return result, errors.New("could not find package to rebuild")
	}

	absPath := filepath.Join(path, packagesPath)
	target, err := repo.FindAndParse(absPath, words[i+1])
	if err != nil {
		return result, err
	}

	dependencies := target.Package.GetDependencies()

	for _, pkg := range packages {
		if pkg.Package.GetName() == target.Package.GetName() {
			return target, nil
		}
		for _, dependency := range dependencies {
			if dependency == pkg.Package.GetName() {
				return target, nil
			}
		}
	}

	return result, errors.New("could not verify the rebuild as needed")

}
