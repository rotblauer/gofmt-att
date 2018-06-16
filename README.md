## `go fmt -w` all the things

You need:
- To look at the `.gofmt-att.json` file. There's a lot of options in there. Pick your poison.
- A Github token with permissions for reading, writing, forking, and PRing. You should have it set as an environment variable. Enter the name of this environment variable in the `"EnvToken": "GITHUB_TOKEN"` field in the bespoke file.
- A git setup on your machine that lets you push to Github with just `git push`. Usernaming-and-passwording isn't currently supported. So figure your ssh shit or whatever out. The program will try to use the `CloneURL` field from the Github API for pushing your changes to your fork.

Run:
- `go run main.go run`
- there's also `go run main.go config` and friends to dump or write a default config


_More about configuration._

- You can use the config file to define _any_ set of commands you want to run on a bunch of repos.
- You can run the spec'd commands on a given repo's whole working directory or file by file (also with white/black lists).
- You can black- and white-list files for git adding or not adding.
- You can spec a lot of other things. If you want to add a different repo provider than Github, use a different database (or "persisting" endpoint in general)... I've tried
to make the architecture extensible.

