package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
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

	// Check if the current run is a PR
	prVal, prExists := os.LookupEnv("GITHUB_EVENT_NAME")
	isPR := prExists && prVal == "pull_request"

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
	fmt.Println()
	fmt.Print(` _           _ _     _  ____             __ _       
| |__  _   _(_) | __| |/ ___|___  _ __  / _(_) __ _ 
| '_ \| | | | | |/ _' | |   / _ \| '_ \| |_| |/ _' |
| |_) | |_| | | | (_| | |__| (_) | | | |  _| | (_| |
|_.__/ \__,_|_|_|\__,_|\____\___/|_| |_|_| |_|\__, |
					       |___/ 
`)
	fmt.Printf("Application Version: v%s\n", config.Global.General.Version)
	fmt.Println()
	fmt.Println()
	fmt.Printf("Containers:\n")
	for container := range containers {
		fmt.Printf("--> %s\n", container)
		output = append(output, container)
	}

	if isPR && config.Global.Git.Token != "" {
		// Grab the current PR number.
		pr, err := repo.PrGetNumber(os.Getenv("GITHUB_REF"))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("[PR %d Detected]\n", pr)
		// Check if packages were changed but resulted in no
		// containers being marked for an update
		if len(containers) == 0 && len(changedPackages) > 0 {
			fmt.Println("Writing Comment for Missing Container Environment...")
			comment := "Although a package changed, no coorisponding " +
				"containers were found to build." +
				" Please make sure to include a `spack.yaml` environment file" +
				" in the `containers/` directory."
			err = repo.PrAddComment(path, config.Global.Git.Token, pr, comment)
			if err != nil {
				log.Fatal(err)
			}
		}
		fmt.Println("Adding Label...")
		if len(containers) == 0 && len(changedPackages) == 0 {
			err = repo.PrAddLabel(path, config.Global.Git.Token, pr, "docs")
		} else {
			err = repo.PrAddLabel(path, config.Global.Git.Token, pr, "build")
		}
		if err != nil {
			log.Fatal(err)
		}
	}
	// Convert results list into JSON
	jsonOutput, _ := json.Marshal(output)
	fmt.Println()
	fmt.Printf("::set-output name=matrix::%s\n", string(jsonOutput))
}
