package engine

import (
	"errors"
	"strings"

	parser "github.com/autamus/binoc/repo"
	"github.com/autamus/buildconfig/repo"
)

func verifyRebuild(packagesPath, commitMsg string, packages []parser.Result) (result parser.Result, err error) {
	name := strings.Fields(commitMsg)[1]
	target, err := repo.FindAndParse(packagesPath, name)
	if err != nil {
		return result, err
	}

	dependencies := target.Package.GetDependencies()

	for _, pkg := range packages {
		for _, dependency := range dependencies {
			if dependency == pkg.Package.GetName() {
				return target, err
			}
		}
	}

	return result, errors.New("could not verify the rebuild as needed")

}
