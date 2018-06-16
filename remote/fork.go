package remote

type ForkConfig struct {
	Org string // organization to fork repo to; if empty will fork to auth'd user
}
