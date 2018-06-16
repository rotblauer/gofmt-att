package remote

import (
	"time"
	"fmt"
	"path/filepath"
	"strings"
)

type RepoT struct {
	Owner    *Owner
	Name     string
	Target   string // filepath of clone
	CloneUrl string
	GitUrl   string
	HTMLUrl  string
}

func (r *RepoT) String() string {
	return fmt.Sprintf("®%s target=%s", r.Ref(), r.Target)
}

func (r *RepoT) Ref() string {
	return filepath.Join(r.Owner.Name, r.Name)
}

func RepoizeRef(fullName string) *RepoT {
	sp := strings.Split(fullName, "/")
	if len(sp) != 2 {
		return nil
	}
	return &RepoT{
		Owner:    &Owner{
			Name:   sp[0],
		},
		Name:     sp[1],
	}
}

func (r *RepoT) GuessURLs(provider string) {
	if provider == "Github" {
		// "html_url": "https://github.com/rotblauer/cathack",
		// "git_url": "git://github.com/rotblauer/cathack.git",
		// "clone_url": "https://github.com/rotblauer/cathack.git",
		r.HTMLUrl = fmt.Sprintf("https://github.com/%s/%s", r.Owner.Name, r.Name)
		r.GitUrl = fmt.Sprintf("git://github.com/%s/%s.git", r.Owner.Name, r.Name)
		r.CloneUrl = fmt.Sprintf("https://github.com/%s/%s.git", r.Owner.Name, r.Name)
	}
}

type RepoListSpec struct {
	*SearchOptions
	OwnerType       string // "user" OR "org" OR "" for either
	Visibility      string // "private", "public", "all"
	Affiliation     string // owner,collaborator,organization_member
	Owner *MatchTextSpec
	Name *MatchTextSpec
	Language        *MatchTextSpec
	Conduct         *MatchTextSpec // eg. "NO ROBOTS" :( or "ROBOTS WECOME"!
	Description     *MatchTextSpec
	IsFork          bool
	IsPrivate       bool
	Archived        bool
	StargazersCount *MatchNSpec
	ForksCount      *MatchNSpec
	WatchersCount   *MatchNSpec
	NetworkCount    *MatchNSpec
	Size            *MatchNSpec
	CreatedAt       *MatchTimeSpec
	UpdatedAt       *MatchTimeSpec
	FmtExpiration  time.Duration `json:"-",toml:"-",yaml:"-"` // parse the below to this on init
	FmtExpiry int      // in DAYS.
}

type SearchOptions struct {
	Sort  string
	Order string
}
//
// {
// "id": 63971126,
// "node_id": "MDEwOlJlcG9zaXRvcnk2Mzk3MTEyNg==",
// "name": "cathack",
// "full_name": "rotblauer/cathack",
// "owner": {
// "login": "rotblauer",
// "id": 20356469,
// "node_id": "MDEyOk9yZ2FuaXphdGlvbjIwMzU2NDY5",
// "avatar_url": "https://avatars3.githubusercontent.com/u/20356469?v=4",
// "gravatar_id": "",
// "url": "https://api.github.com/users/rotblauer",
// "html_url": "https://github.com/rotblauer",
// "followers_url": "https://api.github.com/users/rotblauer/followers",
// "following_url": "https://api.github.com/users/rotblauer/following{/other_user}",
// "gists_url": "https://api.github.com/users/rotblauer/gists{/gist_id}",
// "starred_url": "https://api.github.com/users/rotblauer/starred{/owner}{/repo}",
// "subscriptions_url": "https://api.github.com/users/rotblauer/subscriptions",
// "organizations_url": "https://api.github.com/users/rotblauer/orgs",
// "repos_url": "https://api.github.com/users/rotblauer/repos",
// "events_url": "https://api.github.com/users/rotblauer/events{/privacy}",
// "received_events_url": "https://api.github.com/users/rotblauer/received_events",
// "type": "Organization",
// "site_admin": false
// },
// "private": false,
// "html_url": "https://github.com/rotblauer/cathack",
// "description": "Self-hosted, kitten-friendly collaborative documents written in Golang. Like Google Docs, but eviler.",
// "fork": false,
// "url": "https://api.github.com/repos/rotblauer/cathack",
// "forks_url": "https://api.github.com/repos/rotblauer/cathack/forks",
// "keys_url": "https://api.github.com/repos/rotblauer/cathack/keys{/key_id}",
// "collaborators_url": "https://api.github.com/repos/rotblauer/cathack/collaborators{/collaborator}",
// "teams_url": "https://api.github.com/repos/rotblauer/cathack/teams",
// "hooks_url": "https://api.github.com/repos/rotblauer/cathack/hooks",
// "issue_events_url": "https://api.github.com/repos/rotblauer/cathack/issues/events{/number}",
// "events_url": "https://api.github.com/repos/rotblauer/cathack/events",
// "assignees_url": "https://api.github.com/repos/rotblauer/cathack/assignees{/user}",
// "branches_url": "https://api.github.com/repos/rotblauer/cathack/branches{/branch}",
// "tags_url": "https://api.github.com/repos/rotblauer/cathack/tags",
// "blobs_url": "https://api.github.com/repos/rotblauer/cathack/git/blobs{/sha}",
// "git_tags_url": "https://api.github.com/repos/rotblauer/cathack/git/tags{/sha}",
// "git_refs_url": "https://api.github.com/repos/rotblauer/cathack/git/refs{/sha}",
// "trees_url": "https://api.github.com/repos/rotblauer/cathack/git/trees{/sha}",
// "statuses_url": "https://api.github.com/repos/rotblauer/cathack/statuses/{sha}",
// "languages_url": "https://api.github.com/repos/rotblauer/cathack/languages",
// "stargazers_url": "https://api.github.com/repos/rotblauer/cathack/stargazers",
// "contributors_url": "https://api.github.com/repos/rotblauer/cathack/contributors",
// "subscribers_url": "https://api.github.com/repos/rotblauer/cathack/subscribers",
// "subscription_url": "https://api.github.com/repos/rotblauer/cathack/subscription",
// "commits_url": "https://api.github.com/repos/rotblauer/cathack/commits{/sha}",
// "git_commits_url": "https://api.github.com/repos/rotblauer/cathack/git/commits{/sha}",
// "comments_url": "https://api.github.com/repos/rotblauer/cathack/comments{/number}",
// "issue_comment_url": "https://api.github.com/repos/rotblauer/cathack/issues/comments{/number}",
// "contents_url": "https://api.github.com/repos/rotblauer/cathack/contents/{+path}",
// "compare_url": "https://api.github.com/repos/rotblauer/cathack/compare/{base}...{head}",
// "merges_url": "https://api.github.com/repos/rotblauer/cathack/merges",
// "archive_url": "https://api.github.com/repos/rotblauer/cathack/{archive_format}{/ref}",
// "downloads_url": "https://api.github.com/repos/rotblauer/cathack/downloads",
// "issues_url": "https://api.github.com/repos/rotblauer/cathack/issues{/number}",
// "pulls_url": "https://api.github.com/repos/rotblauer/cathack/pulls{/number}",
// "milestones_url": "https://api.github.com/repos/rotblauer/cathack/milestones{/number}",
// "notifications_url": "https://api.github.com/repos/rotblauer/cathack/notifications{?since,all,participating}",
// "labels_url": "https://api.github.com/repos/rotblauer/cathack/labels{/name}",
// "releases_url": "https://api.github.com/repos/rotblauer/cathack/releases{/id}",
// "deployments_url": "https://api.github.com/repos/rotblauer/cathack/deployments",
// "created_at": "2016-07-22T17:26:08Z",
// "updated_at": "2018-06-13T15:34:44Z",
// "pushed_at": "2018-06-13T15:34:42Z",
// "git_url": "git://github.com/rotblauer/cathack.git",
// "ssh_url": "git@github.com:rotblauer/cathack.git",
// "clone_url": "https://github.com/rotblauer/cathack.git",
// "svn_url": "https://github.com/rotblauer/cathack",
// "homepage": "",
// "size": 53031,
// "stargazers_count": 4,
// "watchers_count": 4,
// "language": "Go",
// "has_issues": true,
// "has_projects": true,
// "has_downloads": true,
// "has_wiki": true,
// "has_pages": false,
// "forks_count": 1,
// "mirror_url": null,
// "archived": false,
// "open_issues_count": 0,
// "license": null,
// "forks": 1,
// "open_issues": 0,
// "watchers": 4,
// "default_branch": "master",
// "permissions": {
// "admin": true,
// "push": true,
// "pull": true
// },
// "allow_squash_merge": true,
// "allow_merge_commit": true,
// "allow_rebase_merge": true,
// "organization": {
// "login": "rotblauer",
// "id": 20356469,
// "node_id": "MDEyOk9yZ2FuaXphdGlvbjIwMzU2NDY5",
// "avatar_url": "https://avatars3.githubusercontent.com/u/20356469?v=4",
// "gravatar_id": "",
// "url": "https://api.github.com/users/rotblauer",
// "html_url": "https://github.com/rotblauer",
// "followers_url": "https://api.github.com/users/rotblauer/followers",
// "following_url": "https://api.github.com/users/rotblauer/following{/other_user}",
// "gists_url": "https://api.github.com/users/rotblauer/gists{/gist_id}",
// "starred_url": "https://api.github.com/users/rotblauer/starred{/owner}{/repo}",
// "subscriptions_url": "https://api.github.com/users/rotblauer/subscriptions",
// "organizations_url": "https://api.github.com/users/rotblauer/orgs",
// "repos_url": "https://api.github.com/users/rotblauer/repos",
// "events_url": "https://api.github.com/users/rotblauer/events{/privacy}",
// "received_events_url": "https://api.github.com/users/rotblauer/received_events",
// "type": "Organization",
// "site_admin": false
// },
// "network_count": 1,
// "subscribers_count": 1
// }
