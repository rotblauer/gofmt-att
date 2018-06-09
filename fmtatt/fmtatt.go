package fmtatt

type FmtAttConfig struct {
	RepoProvider      string
	Identity          AuthIdentity // eg. whilei, ETCDEVTeam, etc.
	ReposFilter       RepoListSpec
	Fmt               FmtConfig
	PullRequestConfig PullRequestConfig
	WalkPattern       WalkPattern
}

type FmtAtt struct {
	Config *FmtAttConfig

	Repoer RepoProvider
	Walker WalkProvider
	Fmters []Fmter
	Giter GitProvider
	Historyers []HistoryProvider

}
