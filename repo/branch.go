package repo

import (
	"errors"
	"log"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// PullBranch attempts to pull the branch from the git origin.
func PullBranch(path string, branchName string) (err error) {
	localBranchReferenceName := plumbing.NewBranchReferenceName(branchName)
	remoteReferenceName := plumbing.NewRemoteReferenceName("origin", branchName)

	r, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	rem, err := r.Remote("origin")
	if err != nil {
		return err
	}

	refs, err := rem.List(&git.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	found := false
	for _, ref := range refs {
		if ref.Name().IsBranch() && ref.Name() == localBranchReferenceName {
			found = true
		}
	}

	if !found {
		return errors.New("branch not found")
	}

	err = r.CreateBranch(&config.Branch{Name: branchName, Remote: "origin", Merge: localBranchReferenceName})
	if err != nil {
		return err
	}
	newReference := plumbing.NewSymbolicReference(localBranchReferenceName, remoteReferenceName)
	err = r.Storer.SetReference(newReference)
	return err
}

// SwitchBranch switches from the current branch to the
// one with the name provided.
func SwitchBranch(path string, branchName string) (err error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	branchRef := plumbing.NewBranchReferenceName(branchName)
	opts := &git.CheckoutOptions{Branch: branchRef}

	err = w.Checkout(opts)
	return err
}

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

// GetOriginalCommit returns the founding commit of the given branch.
func GetOriginalCommit(path string, branchName string, mainName string) (oldCommit *object.Commit, err error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	branchRefName := plumbing.NewBranchReferenceName(branchName)
	mainRefName := plumbing.NewBranchReferenceName(mainName)

	branchRef, err := r.Storer.Reference(branchRefName)
	if err != nil {
		return nil, err
	}

	mainRef, err := r.Storer.Reference(mainRefName)
	if err != nil {
		return nil, err
	}

	branchCommit, err := r.CommitObject(branchRef.Hash())
	if err != nil {
		return nil, err
	}

	mainCommit, err := r.CommitObject(mainRef.Hash())
	if err != nil {
		return nil, err
	}

	gitLog, err := r.Log(&git.LogOptions{})
	if err != nil {
		return oldCommit, err
	}

	err = gitLog.ForEach(func(commit *object.Commit) (err error) {
		isBranchParent, err := commit.IsAncestor(branchCommit)
		if err != nil {
			return err
		}
		isMainParent, err := commit.IsAncestor(mainCommit)
		if err != nil {
			return err
		}
		if isBranchParent && !isMainParent {
			oldCommit, err = commit.Parents().Next()
			return err
		}
		return nil
	})
	if err != nil {
		return oldCommit, err
	}
	if oldCommit == nil {
		return oldCommit, errors.New("no commit found")
	}
	return oldCommit, nil
}
