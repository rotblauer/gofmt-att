package googlegogithub

import (
	"github.com/google/go-github/github"
	"os"
	"github.com/rotblauer/gofmt-att/fmtatt"
	"context"
	"golang.org/x/oauth2"
)

var ghc *github.Client
var ctx context.Context

func NewClient(authID fmtatt.AuthIdentity) *github.Client {
	token := authID.RawToken
	if authID.EnvToken != "" {
		token = os.Getenv(authID.EnvToken)
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken:token},
	)
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, ts)
	ghc = github.NewClient(tc)
	return ghc
}
