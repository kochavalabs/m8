package main

import (
	"log"

	"github.com/kochavalabs/mazzaroth-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
