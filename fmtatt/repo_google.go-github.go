package fmtatt

import (
	"github.com/google/go-github/github"
	"os"
	"github.com/rotblauer/gofmt-att/fmtatt"
	"context"
	"golang.org/x/oauth2"
)

type RepoProvider struct {
	Client *github.Client
	Ctx context.Context
	Username string
}

func NewGoogleGithubProvider(authID fmtatt.AuthIdentity) *RepoProvider {
	token := authID.RawToken
	if authID.EnvToken != "" {
		token = os.Getenv(authID.EnvToken)
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken:token},
	)
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, ts)
	ghc := github.NewClient(tc)

	return &RepoProvider{
		Client: ghc,
		Ctx: ctx,
		Username: authID.Username,
	}
}

func (rp *RepoProvider) GetRepos(reposFilter fmtatt.RepoListSpec) (repos []fmtatt.RepoT, err error) {

	return
}
func (rp *RepoProvider) ForkRepo(rs fmtatt.RepoT) (repo fmtatt.RepoT, err error) {
	return
}
func (rp *RepoProvider) CreatePullRequest(pr fmtatt.SimplePullRequestT) error {
	return nil
}