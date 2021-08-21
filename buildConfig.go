package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	parser "github.com/autamus/binoc/repo"
	"github.com/autamus/buildconfig/config"
	"github.com/autamus/buildconfig/engine"
	"github.com/autamus/buildconfig/repo"
)

type result struct {
	Name      string `json:"name"`
	Arch      string `json:"arch"`
	Container string `json:"container"`
}

func main() {
	// Initialize parser functionality
	parser.Init(strings.Split(config.Global.Parsers.Loaded, ","))

	// Set initial values for Repository
	path := config.Global.Repository.Path
	packagesPath := config.Global.Packages.Path
	containersPath := config.Global.Containers.Path
	defaultEnvPath := config.Global.Containers.DefaultEnvPath
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
		packageContainerDeps, err := repo.IndexPackageContainerDeps(path, defaultEnvPath, containersPath)
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
	buildOutput := make([]result, 0, len(containers))
	pubOutput := make([]result, 0, len(containers))

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
		spackEnv, err := repo.ParseSpackEnv(
			defaultEnvPath,
			filepath.Join(containersPath,
				string(container[0]),
				container,
				"spack.yaml",
			),
		)
		if err == nil {
			pubArches := []string{}
			for _, arch := range spackEnv.Spack.Config.Compiler.Target {
				if arch == "x86_64_v3" {
					buildOutput = append(buildOutput, result{
						Name:      container,
						Arch:      "linux/amd64",
						Container: container,
					})
					pubArches = append(pubArches, "linux/amd64")
				}
				if arch == "aarch64" {
					buildOutput = append(buildOutput, result{
						Name:      container + "-" + "arm",
						Arch:      "linux/arm64",
						Container: container,
					})
					pubArches = append(pubArches, "linux/arm64")
				}
			}
			pubOutput = append(pubOutput, result{
				Name:      container,
				Arch:      strings.Join(pubArches, ","),
				Container: container,
			})
		} else {
			buildOutput = append(buildOutput, result{
				Name:      container,
				Arch:      "linux/amd64",
				Container: container,
			})
			pubOutput = append(pubOutput, result{
				Name:      container,
				Arch:      "linux/amd64",
				Container: container,
			})
		}
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
			comment := "Although a package changed, no corresponding " +
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
	buildJson, _ := json.Marshal(buildOutput)
	pubJson, _ := json.Marshal(pubOutput)
	fmt.Println()
	fmt.Printf("::set-output name=build_matrix::%s\n", string(buildJson))
	fmt.Printf("::set-output name=publish_matrix::%s\n", string(pubJson))
}
