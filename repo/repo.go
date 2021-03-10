package repo

import (
	"bytes"
	"fmt"
	"strings"
)

// GetChangedFiles returns all of the changed files
func GetChangedFiles(path, currentBranch, mainBranch string) (filepaths []string, err error) {
	currentCommit, err := GetCurrentCommit(path)
	if err != nil {
		return filepaths, err
	}
	fmt.Printf("currentCommit: %s\n", currentCommit)

	err = PullBranch(path, mainBranch)
	if err != nil && err.Error() != "branch already exists" {
		return filepaths, err
	}

	originalCommit, err := GetOriginalCommit(path, currentBranch, mainBranch)
	if err != nil {
		return filepaths, err
	}
	fmt.Printf("originalCommit: %s\n", originalCommit)

	diff, err := originalCommit.Patch(currentCommit)
	if err != nil {
		return filepaths, err
	}

	buf := new(bytes.Buffer)
	err = diff.Encode(buf)
	if err != nil {
		return filepaths, err
	}

	lines := strings.Split(buf.String(), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "+++ b/") {
			data := strings.Fields(line)
			filepaths = append(filepaths, strings.TrimPrefix(data[1], "b/"))
		}
	}

	return filepaths, nil
}
