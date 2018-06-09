package cmd

import (
	"github.com/spf13/viper"
)

func viperSetDefaults() {
	viper.Set("repoProvider", "Github")

	viper.Set("identity.username", "")
	viper.Set("identity.rawToken", "")
	viper.Set("identity.envToken", "")

	viper.Set("repos.langauges", []string{"Go"})
	viper.Set("repos.isFork", false)
	viper.Set("repos.sortBy", "updated")
	viper.Set("repos.orderBy", "desc")
	viper.Set("repos.visibility", "visible")

	viper.Set("fmt.commands", []string{
		"gofmt -w",
	})
	viper.Set("fmt.target", ".")
	viper.Set("fmt.perFile", false)
	viper.Set("fmt.files.whitelist", []string{".go$"})
	viper.Set("fmt.files.blacklist", []string{})
	viper.Set("fmt.dirs.whitelist", []string{})
	viper.Set("fmt.dirs.whitelist", []string{"*vendor*", ".git"})

	viper.Set("pr.title", "gofmt: nice and clean. your's truly, the gofmt-att machine")
	viper.Set("pr.head", "")
	viper.Set("pr.base", "master")
	viper.Set("pr.body.text", `
Formattered with :heart: by [gofmt-att](https://github.com/rotblauer/gofmt-att).

> If we got it wrong, or there's a bug or something, please [let us know](https://github.com/rotblauer/gofmt-att/issues/new).
`)
	viper.Set("pr.body.file", "")

	viper.Set("walk.weight.humans", float64(1/2))
	viper.Set("walk.weight.repos", float64(1/2))
	viper.Set("walk.humans.following", float64(1/3))
	viper.Set("walk.humans.followers", float64(1/3))
	viper.Set("walk.humans.orgMembers", float64(1/3))
	viper.Set("walk.repos.starred", float64(1/3))
	viper.Set("walk.repos.forked", float64(1/3))
	viper.Set("walk.repos.authored", float64(1/3))
}
