package repo

import (
	"bytes"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
)

// GetChangedFiles returns all of the changed files
func GetChangedFiles(path, currentBranch, mainBranch string) (filepaths []string, err error) {
	currentCommit, err := GetCurrentCommit(path)
	if err != nil {
		return filepaths, err
	}

	parentCommit, err := currentCommit.Parents().Next()
	if err != nil {
		return filepaths, err
	}

	originalCommit, err := parentCommit.MergeBase(currentCommit)
	if err != nil {
		return filepaths, err
	}

	diff, err := originalCommit[0].Patch(currentCommit)
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

func getURL(path string) (url string, err error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return url, err
	}
	remotes, err := r.Remotes()
	return remotes[0].Config().URLs[0], err
}

func getOwnerName(path string) (repoOwner, repoName string, err error) {
	url, err := getURL(path)
	if err != nil {
		return repoOwner, repoName, err
	}
	repoName = strings.TrimSuffix(filepath.Base(url), filepath.Ext(url))
	repoOwner = filepath.Base(filepath.Dir(url))
	return repoOwner, repoName, nil
}
