package main

import (
	"fmt"
	"os"

	"github.com/garenwen/freebsd-manager/cmd"
)

func main() {
	command := cmd.FreebsdManagerCmd

	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
