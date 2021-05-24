package repo

import (
	"context"
	"path/filepath"
	"strconv"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func PrGetNumber(ref string) (pr int, err error) {
	prString := filepath.Base(filepath.Dir(ref))
	return strconv.Atoi(prString)
}

func PrAddComment(path, gitToken string, pr int, comment string) (err error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gitToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	repoOwner, repoName, err := getOwnerName(path)
	if err != nil {
		return err
	}

	_, _, err = client.PullRequests.CreateComment(
		ctx,
		repoOwner,
		repoName,
		pr,
		&github.PullRequestComment{
			Body: &comment,
		})

	return err
}

func PrAddLabel(path, gitToken string, pr int, label string) (err error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gitToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	repoOwner, repoName, err := getOwnerName(path)
	if err != nil {
		return err
	}

	_, _, err = client.Issues.AddLabelsToIssue(
		ctx,
		repoOwner,
		repoName,
		pr,
		[]string{label},
	)

	return err
}
