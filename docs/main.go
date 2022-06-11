package main

import (
	"github.com/rafaelcalleja/go-bpkg/pkg/cmd"
	"github.com/spf13/cobra/doc"
	"log"
)

func main() {
	rootCmd := cmd.Main()

	err := doc.GenMarkdownTree(rootCmd, "./docs/")
	if err != nil {
		log.Fatal(err)
	}
}
