package fmtatt

// WalkProvider implements the logic required to select a next batch of repos to fetch.
type WalkProvider interface {
	StepNext(pattern WalkPattern, history HistoryProvider) (rs RepoListSpec, err error)
}

// WalkPattern decides how to select a next batch of repos of solve.
// HumansWeight + ReposWeight must equal 1.
type WalkPattern struct {
	// Preference for using human connections to prioritize next repos.
	HumansWeight float64
	// Preference for using repo connections to prioritize next repos.
	ReposWeight float64

	// Distance is how far from the catalyst to venture.
	// 0 for stop after first batch of repos.
	// -1 for the whole world.
	// 6 for six degrees of separation.
	Distance int
	WalkHumansPattern
	WalkReposPatterns
}

// WalkHumansPattern decides how to balance human connection variables when choosing
// a next batch of repos to solve.
type WalkHumansPattern struct {
	Following bool
	Followers bool
	OrgMembers bool
}

// WalkReposPatterns decides how to balance repo list sources.
type WalkReposPatterns struct {
	Starred bool
	Forked bool
	Authored bool
}
