package remote

import (
	"github.com/google/go-github/github"
	"regexp"
)

func (rp *GoogleGithubRepoProvider) FilterRepo(rawRepo interface{}, spec *RepoListSpec) (err *ErrFilteredT) {
	repo := (rawRepo).(*github.Repository)
	if err = filterTextSpec("repo", spec.Owner, repo.GetOwner().GetLogin()); err != nil {
		return
	}
	if err = filterTextSpec("repo", spec.Name, repo.GetName()); err != nil {
		return
	}
	if err = filterTextSpec("repo", spec.Description, repo.GetDescription()); err != nil {
		return
	}
	if err = filterTextSpec("repo", spec.Language, repo.GetLanguage()); err != nil {
		return
	}
	if err = filterTextSpec("repo", spec.Conduct, repo.GetCodeOfConduct().GetBody()); err != nil {
		return
	}
	if err = filterNSpec("repo", spec.StargazersCount, repo.GetStargazersCount()); err != nil {
		return
	}
	if err = filterNSpec("repo", spec.ForksCount, repo.GetForksCount()); err != nil {
		return
	}
	if err = filterNSpec("repo", spec.WatchersCount, repo.GetWatchersCount()); err != nil {
		return
	}
	if err = filterNSpec("repo", spec.NetworkCount, repo.GetNetworkCount()); err != nil {
		return
	}
	if err = filterNSpec("repo", spec.Size, repo.GetSize()); err != nil {
		return
	}

	if err = filterTimeSpec("repo", spec.CreatedAt, repo.GetCreatedAt().Time); err != nil {
		return
	}
	if err = filterTimeSpec("repo", spec.UpdatedAt, repo.GetUpdatedAt().Time); err != nil {
		return
	}
	if ok := filterBoolSpec("repo", spec.IsFork, repo.GetFork()); !ok {
		err = &ErrFilteredT{
			Reason:   "isFork",
			Value:    repo.GetFork(),
			Resource: "repo",
		}
		return
	}
	if ok := filterBoolSpec("repo", spec.IsPrivate, repo.GetPrivate()); !ok {
		err = &ErrFilteredT{
			Reason:   "isPrivate",
			Value:    repo.GetPrivate(),
			Resource: "repo",
		}
		return
	}
	return
}

func (rp *GoogleGithubRepoProvider) FilterOwner(rawOwner interface{}, spec *OwnerListSpec) (err *ErrFilteredT) {
	owner := (rawOwner).(*github.User)
	wantName := rp.ownerSpecs.Owner.Name
	if wantName != "" {
		re := regexp.MustCompile(wantName)
		if !re.MatchString(owner.GetLogin()) {
			err = &ErrFilteredT{
				Reason:   "name",
				Value:    owner.Name,
				Resource: "user",
			}
			return
		}
	}

	if rp.ownerSpecs.Hireable != "" {
		gotHirable := owner.GetHireable()
		wantHirable := rp.ownerSpecs.Hireable == "true"
		if gotHirable != wantHirable {
			err = &ErrFilteredT{
				Reason:   "hireable",
				Value:    gotHirable,
				Resource: "user",
			}
			return
		}
	}

	if err = filterNSpec("user", rp.ownerSpecs.PublicGistsN, owner.GetPublicGists()); err != nil {
		return
	}
	if err = filterNSpec("user", rp.ownerSpecs.PublicReposN, owner.GetPublicRepos()); err != nil {
		return
	}
	if err = filterNSpec("user", rp.ownerSpecs.FollowersN, owner.GetFollowers()); err != nil {
		return
	}
	if err = filterNSpec("user", rp.ownerSpecs.FollowingN, owner.GetFollowing()); err != nil {
		return
	}
	if err = filterNSpec("user", rp.ownerSpecs.CollaboratorsN, owner.GetCollaborators()); err != nil {
		return
	}
	return
}
