package remote

import (
	"fmt"
)

type Owner struct {
	Name   string
	KindOf string // "Owner" "Organization"
}

func (o *Owner) String() string {
	if o.KindOf == "Organization" {
		return fmt.Sprintf("👪  %s", o.Name)
	}
	return fmt.Sprintf("🤓  %s", o.Name)
}

// OwnerListSpec can be user or org
type OwnerListSpec struct {
	Owner
	FollowingN     *MatchNSpec
	FollowersN     *MatchNSpec
	PublicReposN   *MatchNSpec
	PublicGistsN   *MatchNSpec
	CollaboratorsN *MatchNSpec
	Hireable       string // "true", "false", ""
}

type UserListSpec struct {
	OwnerListSpec
}

type OrganizationListSpec struct {
	OwnerListSpec
}

// {
// "login": "whilei",
// "id": 10228550,
// "node_id": "MDQ6VXNlcjEwMjI4NTUw",
// "avatar_url": "https://avatars1.githubusercontent.com/u/10228550?v=4",
// "gravatar_id": "",
// "url": "https://api.github.com/users/whilei",
// "html_url": "https://github.com/whilei",
// "followers_url": "https://api.github.com/users/whilei/followers",
// "following_url": "https://api.github.com/users/whilei/following{/other_user}",
// "gists_url": "https://api.github.com/users/whilei/gists{/gist_id}",
// "starred_url": "https://api.github.com/users/whilei/starred{/owner}{/repo}",
// "subscriptions_url": "https://api.github.com/users/whilei/subscriptions",
// "organizations_url": "https://api.github.com/users/whilei/orgs",
// "repos_url": "https://api.github.com/users/whilei/repos",
// "events_url": "https://api.github.com/users/whilei/events{/privacy}",
// "received_events_url": "https://api.github.com/users/whilei/received_events",
// "type": "Owner",
// "site_admin": false,
// "name": "ia",
// "company": null,
// "blog": "",
// "location": null,
// "email": null,
// "hireable": null,
// "bio": "Creative writer at @ETCDEVTeam.",
// "public_repos": 138,
// "public_gists": 60,
// "followers": 45,
// "following": 140,
// "created_at": "2014-12-18T04:37:38Z",
// "updated_at": "2018-06-15T05:36:20Z"
// }

//
// {
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
// "site_admin": false,
// "name": "rotblauer",
// "company": null,
// "blog": "http://www.rotblauer.com",
// "location": "Lewes, England",
// "email": "rotblauer@gmail.com",
// "hireable": null,
// "bio": null,
// "public_repos": 80,
// "public_gists": 0,
// "followers": 0,
// "following": 0,
// "created_at": "2016-07-08T12:45:03Z",
// "updated_at": "2017-02-13T23:36:12Z"
// }
