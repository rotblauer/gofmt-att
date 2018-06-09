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

	oneHalf := float64(1)/float64(2)
	oneThird := float64(1)/float64(3)
	viper.Set("walk.weight.humans", oneHalf)
	viper.Set("walk.weight.repos", oneHalf)
	viper.Set("walk.humans.following", oneThird)
	viper.Set("walk.humans.followers", oneThird)
	viper.Set("walk.humans.orgMembers", oneThird)
	viper.Set("walk.repos.starred", oneThird)
	viper.Set("walk.repos.forked", oneThird)
	viper.Set("walk.repos.authored", oneThird)
}
