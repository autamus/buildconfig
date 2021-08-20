package repo

import (
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type SpackEnv struct {
	Spack Spack `yaml:"spack"`
}

type Spack struct {
	Specs     []string                 `yaml:"specs,omitempty"`
	View      bool                     `yaml:"view"`
	Packages  map[string]SpackPackages `yaml:"packages,omitempty"`
	Config    SpackConfig              `yaml:"config"`
	Container SpackContainer           `yaml:"container"`
	Mirrors   map[string]string        `yaml:"mirrors,omitempty"`
}

type SpackPackages struct {
	Target []string `yaml:"target,omitempty"`
}

type SpackConfig struct {
	Concretizer             string                 `yaml:"concretizer,omitempty"`
	Compiler                SpackConfigCompiler    `yaml:"compiler,omitempty"`
	InstallMissingCompilers bool                   `yaml:"install_missing_compilers"`
	InstallTree             SpackConfigInstallTree `yaml:"install_tree,omitempty"`
}

type SpackConfigInstallTree struct {
	Root         string `yaml:"root,omitempty"`
	PaddedLength int    `yaml:"padded_length,omitempty"`
}

type SpackConfigCompiler struct {
	Target []string `yaml:"target,omitempty"`
}

type SpackContainer struct {
	OSPackages SpackContainerPackages `yaml:"os_packages,omitempty"`
	Images     SpackContainerImages   `yaml:"images,omitempty"`
	Strip      bool                   `yaml:"strip"`
}

type SpackContainerImages struct {
	Build string `yaml:"build,omitempty"`
	Final string `yaml:"final,omitempty"`
}

type SpackContainerPackages struct {
	Build []string `yaml:"build,omitempty"`
	Final []string `yaml:"final,omitempty"`
}

func defaultEnv(defaultPath string) (output SpackEnv, err error) {
	input, err := ioutil.ReadFile(defaultPath)
	if err != nil {
		return output, err
	}
	err = yaml.Unmarshal(input, &output)
	return output, err
}

// ParseSpackEnv parses a spack environment into a go struct.
func ParseSpackEnv(defaultPath, containerPath string) (result SpackEnv, err error) {
	result, err = defaultEnv(defaultPath)
	if err != nil {
		return
	}
	input, err := ioutil.ReadFile(containerPath)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(input, &result)
	if err != nil {
		return
	}

	// Clean specs from variant/version information.
	specs := []string{}
	for _, spec := range result.Spack.Specs {
		i := strings.IndexFunc(spec, versend)
		if i > 0 {
			specs = append(specs, spec[:i])
		} else {
			specs = append(specs, spec)
		}
	}
	result.Spack.Specs = specs

	return result, nil
}

// versend returns true at the end of the name of a dependency
func versend(input rune) bool {
	for _, c := range []rune{'@', '~', '+', '%'} {
		if input == c {
			return true
		}
	}
	return false
}
