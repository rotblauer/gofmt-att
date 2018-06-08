package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"net/http"
	"gopkg.in/src-d/go-git.v4"
	gitConfig "gopkg.in/src-d/go-git.v4/config"
	gitTransport "gopkg.in/src-d/go-git.v4/plumbing/transport"
	gitObject "gopkg.in/src-d/go-git.v4/plumbing/object"
	"time"
	"os/exec"
)

func fixemloop(c *http.CLient) {
	gitTransport.AuthMethod()
}

func main() {

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// list all repositories for a user
	repos, _, err := client.Repositories.List(ctx, "whilei", nil)
	if err != nil {
		log.Fatalln(err)
	}
	for _, repo := range repos {
		if !repo.GetFork() && repo.GetLanguage() == "Go" {
			fmt.Println(repo.GetCloneURL())
		}
	}

	// client.Octocat()

	users, _, err := client.Users.ListFollowing(ctx, "whilei", nil)
	if err != nil {
		log.Fatalln(err)
	}
	for _, user := range users {
		name := user.Name
		id := user.ID
	}

	starredRepos, _, err := client.Activity.ListStarred(ctx, "whilei", nil)
	if err != nil {
		log.Fatalln(err)
	}
	for _, starredRepo := range starredRepos {
		starredRepo.Repository.GetLanguage()
		starredRepo.Repository.GetCloneURL()

	}

	_, err := git.PlainClone("/tmp/foo", false, &git.CloneOptions{
		URL:      "https://github.com/src-d/go-git",
		Progress: os.Stdout,
	})

	forkedRepo, _, err := client.Repositories.CreateFork(ctx, "whilei", "fillin", nil)
	if err != nil {
		log.Fatalln(err)
	}

	r, err := git.PlainOpen("path/to/repo")
	r.CreateBranch(&gitConfig.Branch{
		Name: "gofmt-att",
		Remote: "origin",
		Merge: "refs/heads/gofmt-att",
	})
	w, err := r.Worktree()
	if err != nil {
		// handle error
	}
	_, err = w.Add(".")
	if err != nil {
		// handle error
	}
	status, err := w.Status()
	if err != nil {
		// handle error
	}
	status.IsClean()
	commit, err := w.Commit("example go-git commit", &git.CommitOptions{
		Author: &gitObject.Signature{
			Name:  "John Doe",
			Email: "john@doe.org",
			When:  time.Now(),
		},
	})

	if err := r.Push(&git.PushOptions{}); err != nil {
		// handle error
	}

	exec.Command("gofmt ...")
	// if err := r.Push(&git.PushOptions{
	// 	// RemoteName is the name of the remote to be pushed to.
	// 	RemoteName: "origin",
	// 	// RefSpecs specify what destination ref to update with what source
	// 	// object. A refspec with empty src can be used to delete a reference.
	// 	RefSpecs: []gitConfig.RefSpec{"+refs/heads/*:refs/remotes/origin/*"},
	// 	// Auth credentials, if required, to use with the remote repository.
	// 	Auth transport.AuthMethod
	// 	// Progress is where the human readable information sent by the server is
	// 	// stored, if nil nothing is stored.
	// 	Progress sideband.Progress
	// }); err != nil {
	// 	// handle error
	// }

	newPR := &github.NewPullRequest{
		Title:               github.String("My awesome pull request"),
		Head:                github.String("branch_to_merge"),
		Base:                github.String("master"),
		Body:                github.String("This is the description of the PR created with the package `github.com/google/go-github/github`"),
		MaintainerCanModify: github.Bool(true),
	}

	pr, _, err := client.PullRequests.Create(context.Background(), "myOrganization", "myRepository", newPR)
	if err != nil {
		fmt.Println(err)
		return
	}

}
