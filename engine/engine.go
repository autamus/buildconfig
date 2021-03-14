package engine

import (
	"strings"

	parser "github.com/autamus/binoc/repo"
	"github.com/go-git/go-git/v5"
)

// FindTarget attempts to find the target of a build.
func FindTarget(path, packagesPath string, packages []parser.Result) (result parser.Result, err error) {
	commitMsg, err := getCommitMsg(path)
	if err != nil {
		return result, err
	}

	commitMsg = strings.ToLower(commitMsg)

	switch {
	case strings.Contains(commitMsg, "rebuild"):
		result, err = verifyRebuild(path, packagesPath, commitMsg, packages)
		if err != nil {
			return result, err
		}
		break

	default:
		result, err = findUpdate(packagesPath, commitMsg, packages)
		if err != nil {
			return result, err
		}
	}

	return result, nil
}

func getCommitMsg(path string) (result string, err error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return result, err
	}

	ref, err := r.Head()
	if err != nil {
		return result, err
	}

	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return result, err
	}

	return commit.Message, nil
}
