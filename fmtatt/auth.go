package fmtatt

import (
	"os"
	"fmt"
)

func mustGetTokenFromConfig(c *Config) (token string) {
	token = c.Identity.RawToken
	if c.Identity.EnvToken != "" {
		token = os.Getenv(c.Identity.EnvToken)
	}
	if token == "" {
		fmt.Println("missing token")
		os.Exit(1)
	}
	return
}
