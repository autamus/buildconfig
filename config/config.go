package config

import (
	"os"
	"reflect"
	"strings"
)

// Config defines the configuration struct for importing settings from ENV Variables
type Config struct {
	General    general
	Packages   packages
	Containers containers
	Repository repository
	Parsers    parsers
}

type general struct {
	Version string
}

type packages struct {
	Path string
}

type containers struct {
	Path string
}

type repository struct {
	Path          string
	DefaultBranch string
}

type parsers struct {
	Loaded string
}

var (
	// Global is the configuration struct for the application.
	Global Config
)

func init() {
	defaultConfig()
	parseConfigEnv()
}

func defaultConfig() {
	Global.General.Version = "0.0.1"
	Global.Containers.Path = "containers"
	Global.Packages.Path = "spack/"
	Global.Repository.Path = "."
	Global.Repository.DefaultBranch = "main"
	Global.Parsers.Loaded = "spack"
}

func parseConfigEnv() {
	numSubStructs := reflect.ValueOf(&Global).Elem().NumField()
	for i := 0; i < numSubStructs; i++ {
		iter := reflect.ValueOf(&Global).Elem().Field(i)
		subStruct := strings.ToUpper(iter.Type().Name())

		structType := iter.Type()
		for j := 0; j < iter.NumField(); j++ {
			fieldVal := iter.Field(j).String()
			if fieldVal != "Version" {
				fieldName := structType.Field(j).Name
				for _, prefix := range []string{"BC", "INPUT"} {
					evName := prefix + "_" + subStruct + "_" + strings.ToUpper(fieldName)
					evVal, evExists := os.LookupEnv(evName)
					if evExists && evVal != fieldVal {
						iter.FieldByName(fieldName).SetString(evVal)
					}
				}
			}
		}
	}
}
