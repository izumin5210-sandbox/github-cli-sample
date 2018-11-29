package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/go-github/v19/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"

	"github.com/izumin5210-sandbox/github-cli-sample/pkg/ghcp"
)

func NewGhcpCommand(c *ghcp.Ctx) *cobra.Command {
	cmd := &cobra.Command{
		Use: "ghcp",
		RunE: func(_ *cobra.Command, args []string) error {
			owner, repo, branch := args[0], args[1], args[2]

			ctx := context.Background()
			gc := githubClient(ctx)

			// get master's ref
			masterRef, _, err := gc.Git.GetRef(ctx, owner, repo, "heads/master")
			if err != nil {
				return err
			}

			// create branch
			featureRef, _, err := gc.Git.CreateRef(ctx, owner, repo, &github.Reference{
				Ref:    "heads/" + branch,
				Object: masterRef.Object,
			})
			if err != nil {
				return err
			}

			// get HEAD commit
			headCommit, _, err := gc.Git.GetCommit(ctx, owner, repo, featureRef.Object.SHA)
			if err != nil {
				return err
			}

			// create blob
			blob, _, err := gc.Git.CreateBlob(ctx, owner, repo, &github.Blob{
				Content:  "foobarbaz",
				Encoding: "utf-8",
			})
			if err != nil {
				return err
			}

			tree, _, err := gc.Git.CreateTree(ctx, owner, repo, headCommit.Tree.SHA, []github.TreeEntry{
				{
					Path: "test.txt",
					Mode: "100644",
					Type: "blob",
					SHA:  blob.SHA,
				},
			})
			if err != nil {
				return err
			}

			commit, _, err := gc.Git.CreateCommit(ctx, owner, repo, &github.Commit{
				Message: "test commit",
				Author: &github.CommitAuthor{
					Name:  "izumin5210",
					Email: "m@izum.in",
					Date:  time.Now().Format("2006-01-02T15:04:05-0700"),
				},
				Parents: []github.Commit{headCommit},
				Tree:    tree.SHA,
			})
			if err != nil {
				return err
			}

			_, _, err = gc.Git.UpdateRef(ctx, owner, repo, &github.Reference{
				Ref: branch,
				SHA: commit.SHA,
			}, false)
			if err != nil {
				return err
			}

			pr, _, err := gc.PullRequests.Create(ctx, owner, repo, &github.NewPullRequest{
				Title: "test pull request",
				Head:  branch,
				Base:  "master",
				Body:  "This is a test pull request.",
			})
			if err != nil {
				return err
			}

			fmt.Fprintln(c.IO.Out, pr.URL)

			return nil
		},
	}

	return cmd
}

func githubClient(ctx context.Context) *github.Client {
	hc := httpClient(ctx)
	gc := github.NewClient(hc)
	return gc
}

func httpClient(ctx context.Context) *http.Client {
	token := os.Getenv("GITHUB_TOKEN")
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	return tc
}
