package main

import (
	"blockchain/cli"
	"os"
)

func main() {
	// help to properly close database
	defer os.Exit(0)
	// init cli
	cli := cli.CommandLine{}
	// run cli
	cli.Run()
}
