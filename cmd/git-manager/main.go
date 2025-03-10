package main

import (
	"fmt"
	"os"

	"github.com/ingshtrom/git-manager/cmd/git-manager/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
