package main

import (
	"encoding/json"
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
	containersPath := config.Global.Containers.Path
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

	// Get a list of all of the containers directly changed in the commit.
	containers, err := repo.GetChangedContainers(path, containersPath, filepaths)
	if err != nil {
		log.Fatal(err)
	}

	// Get a list of all of the packages modified in the commit.
	changedPackages, err := repo.GetChangedPackages(path, packagesPath, filepaths)
	if err != nil {
		log.Fatal(err)
	}

	if len(changedPackages) > 0 {
		// Build a map of packages (values) that rely on a package (key).
		reversePackageDeps, err := repo.IndexReverseDependencies(path, packagesPath)
		if err != nil {
			log.Fatal(err)
		}
		// Build a map of containers (values) that rely on a package (key).
		packageContainerDeps, err := repo.IndexPackageContainerDeps(path, containersPath)
		if err != nil {
			log.Fatal(err)
		}

		packages := engine.GetAllPackageBuilds(changedPackages, reversePackageDeps)
		for app := range packages {
			for _, container := range packageContainerDeps[app] {
				containers[container] = true
			}
		}
	}

	// Initialize list for keys of containers
	output := make([]string, 0, len(containers))

	// Print BuildConfig Report
	fmt.Println("[BuildConfig]")
	fmt.Printf("v%s\n", config.Global.General.Version)
	fmt.Println()
	fmt.Printf("Containers:\n")
	for container := range containers {
		fmt.Printf("--> %s\n", container)
		output = append(output, container)
	}
	// Convert results list into JSON
	jsonOutput, _ := json.Marshal(output)
	fmt.Println()
	fmt.Printf("::set-output name=matrix::%s\n", string(jsonOutput))
}
