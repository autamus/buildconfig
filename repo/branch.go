package repo

import (
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// GetBranchName returns the name of the current branch.
func GetBranchName(path string) (name string, err error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return name, err
	}
	h, err := r.Head()
	if err != nil {
		return name, err
	}
	name = strings.TrimPrefix(h.Name().String(), "refs/heads/")

	return name, nil
}

// GetCurrentCommit returns the current head commit from a branch.
func GetCurrentCommit(path string) (commit *object.Commit, err error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return commit, err
	}

	ref, err := r.Head()
	if err != nil {
		return commit, err
	}

	commit, err = r.CommitObject(ref.Hash())
	if err != nil {
		return commit, err
	}

	return commit, nil
}
