package main

import (
	"os"

	"github.com/Jdubrick/rhdh-profile/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
