package fmtatt

import (
	"fmt"
	"log"
	"os"

	"io/ioutil"
	"flag"

	"github.com/google/go-github/github"
	"context"
)

var (
	// This needs to be set in order to use authenticated Github API requests.
	GITHUB_TOKEN_TEST = "GITHUB_TOKEN"

	workdir_test = os.TempDir()
)

func init() {
	d, err := ioutil.TempDir("", "gofmt-att")
	if err != nil {
		log.Fatalln(err)
	}
	workdir_test = d

	flag.StringVar(&GITHUB_TOKEN_TEST, "github-token", "GITHUB_TOKEN", "Github token to use for API authentication. Can be EITHER the name of an environment variable, or a raw token.")
	flag.Parse()
}

func getAllUserRepos(ctx context.Context, c *github.Client, user string) (repos []*github.Repository, err error) {
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	for {
		pagedRepos, resp, err := c.Repositories.List(ctx, user, opt)
		if err != nil {
			return nil, err
		}
		repos = append(repos, pagedRepos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return
}

// func getAllStarredGoRepos(user string) {
//
// }

func main() {
	ctx, client := NewGithubClientAndContext(GITHUB_TOKEN_TEST)

	// list all repositories for a user
	// repos, _, err := client.Repositories.List(ctx, "whilei", nil)
	repos, err := getAllUserRepos(ctx, client, "whilei")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("- Got %d Golang repos\n", len(repos))
	for _, repo := range repos {
		if !repo.GetFork() && repo.GetLanguage() == "Go" {
			fmt.Printf("lang: %s clone_url: %s git_url: %s\n", repo.GetLanguage(), repo.GetCloneURL(), repo.GetGitURL())
		}
	}

	// for _, repoStarred := range repos {
	// 	repo := repoStarred.Repository
	// 	if !repo.GetFork() && repo.GetLanguage() == "Go" {
	// 		fmt.Printf("lang: %s clone_url: %s git_url: %s\n", repo.GetLanguage(), repo.GetCloneURL(), repo.GetGitURL())
	// 	}
	// }

}
