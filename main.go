package main

import (
	"os"

	"github.com/kochavalabs/mazzaroth-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(0)
	}
}
