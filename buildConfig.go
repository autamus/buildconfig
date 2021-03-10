package main

import (
	"fmt"
	"log"
	"strings"

	parser "github.com/autamus/binoc/repo"
	"github.com/autamus/buildconfig/config"
	"github.com/autamus/buildconfig/engine"
	"github.com/autamus/buildconfig/repo"
)

func main() {
	// Initialize parser functionality
	parser.Init(strings.Split(config.Global.Parsers.Loaded, ","))

	// Set inital values for Repository
	path := config.Global.Repository.Path
	packagesPath := config.Global.Packages.Path
	mainBranch := config.Global.Repository.DefaultBranch

	currentBranch, err := repo.GetBranchName(path)
	if err != nil {
		log.Fatal(err)
	}

	// Get a list of all of the changed files in the commit.
	filepaths, err := repo.GetChangedFiles(path, currentBranch, mainBranch)
	if err != nil {
		log.Fatal(err)
	}

	// Get a list of all of the packages in the commit.
	packages, err := repo.GetPackages(path, filepaths)
	if err != nil {
		log.Fatal(err)
	}

	result, err := engine.FindTarget(path, packagesPath, packages)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", result.Package.GetName())

	version := strings.Join(result.Package.GetLatestVersion(), ".")

	fmt.Printf("%s\n", version)

}
