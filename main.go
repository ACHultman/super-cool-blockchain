package main

import (
	"blockchain/wallet"
	"os"
)

func main() {
	// help to properly close database
	defer os.Exit(0)
	// init cli
	/*	cli := cli.CommandLine{}
		// run cli
		cli.Run()*/
	w := wallet.MakeWallet()
	w.Address()
}
