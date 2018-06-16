## `gofmt -w` all the things

This is a program that walks around Github, grabs Go repos, clones 'em, and runs `gofmt -w .` on 'em.

If after running the notorious gofmt, changes have been (eg. `git status` is dirty), then it adds, commits,
forks, pushes, and finally puts a PR up. 

If you're one of the maintainers or owners of a project and I've bothered you or messed up in my proposed change,
please file an issue so this work can be as helpful and positive as possible. 

I've been running with some configured blacklisted regex's (like vendor/ dirs) to exclude from the `git add` portion to try to keep changes pertinent.

Want to run it yourself? Cool!

To run it, you need:
- To look at the `.gofmt-att.json` file. There's a lot of options in there. Pick your poison.
- A Github token with permissions for reading, writing, forking, and PRing. You should have it set as an environment variable. Enter the name of this environment variable in the `"EnvToken": "GITHUB_TOKEN"` field in the bespoke file.
- A git setup on your machine that lets you push to Github with just `git push`. Usernaming-and-passwording isn't currently supported. So figure your ssh shit or whatever out. The program will try to use the `CloneURL` field from the Github API for pushing your changes to your fork.

Run:
- `go run main.go run`
- there's also `go run main.go config` and friends to dump or write a default config

_More about configuration and other stuff._

- You can use the config file to define _any_ set of commands you want to run on a bunch of repos. So if you can find a spellcheck program that actually works...
- You can run the spec'd commands on a given repo's whole working directory or file by file (also with white/black lists).
- You can black- and white-list files for git adding or not adding.
- You can spec a lot of other things. If you want to add a different repo provider than Github, use a different database (or "persisting" endpoint in general)... I've tried
to make the architecture extensible.

> Rob Pike touching on gofmt: https://www.youtube.com/watch?v=PAAkCSZUG1c&t=8m43s