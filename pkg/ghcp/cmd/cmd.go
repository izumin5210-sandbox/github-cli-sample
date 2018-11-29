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

			svc := &service{
				Client: gc,
				Author: Author{
					Name:  "izumin5210",
					Email: "m@izum.in",
				},
			}

			ref, err := svc.FindOrCreateBranch(ctx, owner, repo, "master", branch)
			if err != nil {
				return err
			}

			err = svc.CreateOrUpdateCommit(ctx, owner, repo, ref, []File{
				{
					Path:    "test.txt",
					Content: "foobarbaz",
				},
			}, "test commit")
			if err != nil {
				return err
			}

			pr, err := svc.CreateOrUpdatePullRequest(ctx, owner, repo, "master", branch, "test pull request", "This is test pull request")
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

type service struct {
	*github.Client

	Author Author
}

type Author struct {
	Name  string
	Email string
}

type File struct {
	Path    string
	Content string
}

func (s *service) FindOrCreateBranch(ctx context.Context, owner, repo string, baseBranch, featureBranch string) (
	*github.Reference,
	error,
) {
	baseBranch = "heads/" + baseBranch
	featureBranch = "heads/" + featureBranch

	// TODO: find an existing branch

	// get master's ref
	masterRef, _, err := s.Git.GetRef(ctx, owner, repo, baseBranch)
	if err != nil {
		return nil, err
	}

	// create branch
	featureRef, _, err := s.Git.CreateRef(ctx, owner, repo, &github.Reference{
		Ref:    &featureBranch,
		Object: masterRef.Object,
	})
	if err != nil {
		return nil, err
	}

	return featureRef, nil
}

func (s *service) CreateOrUpdateCommit(ctx context.Context, owner, repo string, ref *github.Reference, files []File, msg string) error {
	headCommit, _, err := s.Git.GetCommit(ctx, owner, repo, *ref.Object.SHA)
	if err != nil {
		return err
	}

	var (
		encoding  = "utf-8"
		entryMode = "100644"
		entryType = "blob"
		entries   []github.TreeEntry
	)
	for _, f := range files {
		blob, _, err := s.Git.CreateBlob(ctx, owner, repo, &github.Blob{
			Content:  &f.Content,
			Encoding: &encoding,
		})
		if err != nil {
			return err
		}

		entries = append(entries, github.TreeEntry{
			Path: &f.Path,
			Mode: &entryMode,
			Type: &entryType,
			SHA:  blob.SHA,
		})
	}

	tree, _, err := s.Git.CreateTree(ctx, owner, repo, *headCommit.Tree.SHA, entries)
	if err != nil {
		return err
	}

	date := time.Now()
	commit, _, err := s.Git.CreateCommit(ctx, owner, repo, &github.Commit{
		Message: &msg,
		Author: &github.CommitAuthor{
			Name:  &s.Author.Name,
			Email: &s.Author.Email,
			Date:  &date,
		},
		Parents: []github.Commit{*headCommit},
		Tree:    tree,
	})
	if err != nil {
		return err
	}

	_, _, err = s.Git.UpdateRef(ctx, owner, repo, &github.Reference{
		Ref: ref.Ref,
		Object: &github.GitObject{
			SHA: commit.SHA,
		},
	}, false)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) CreateOrUpdatePullRequest(ctx context.Context, owner, repo string, baseBranch, featureBranch, title, body string) (*github.PullRequest, error) {
	pr, _, err := s.PullRequests.Create(ctx, owner, repo, &github.NewPullRequest{
		Title: &title,
		Head:  &featureBranch,
		Base:  &baseBranch,
		Body:  &body,
	})
	return pr, err
}
