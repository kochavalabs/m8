package main

import (
	"os"

	"github.com/kochavalabs/m8/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(0)
	}
}
