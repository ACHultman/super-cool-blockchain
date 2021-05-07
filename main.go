package main

import (
	"blockchain/blockchain"
	"flag"
	"fmt"
	"github.com/dgraph-io/badger"
	"os"
	"runtime"
	"strconv"
)

type CommandLine struct {
	blockChain *blockchain.BlockChain
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" add -block BLOCK_DATA - add a block to the chain")
	fmt.Println(" print - Prints the blocks in the chain")
}

func (cli *CommandLine) validateArgs() {
	// check number of args user has entered
	if len(os.Args) < 2 {
		// if no args
		// print usage
		cli.printUsage()
		// exit app by shutting down Go routine (to allow badger garbage collection)
		runtime.Goexit()
	}
}

func (cli *CommandLine) addBlock(data string) {
	cli.blockChain.AddBlock(data)
	fmt.Println("Added block!")
}

func (cli *CommandLine) printChain() {
	iter := cli.blockChain.Iterator()

	for {
		block := iter.Next()
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)

		// run proof of work
		pow := blockchain.NewProof(block)
		// print validation
		fmt.Printf("POW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PRevHash) == 0 {
			// if on Genesis block
			break
		}
	}
}

func (cli *CommandLine) run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printchainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	addBlockData := addBlockCmd.String("block", "", "Block data")

	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		blockchain.HandleError(err)
	case "print":
		err := printchainCmd.Parse(os.Args[2:])
		blockchain.HandleError(err)
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if addBlockCmd.Parsed() {
		// if add block command parsed
		if *addBlockData == "" {
			// if block data empty
			// print add block command usage
			addBlockCmd.Usage()
			// go exit
			runtime.Goexit()
		}
		// block data received
		// add block with data
		cli.addBlock(*addBlockData)
	}

	if printchainCmd.Parsed() {
		// if print chain command parsed
		cli.printChain()
	}
}

func main() {
	// help to properly close database
	defer os.Exit(0)
	// init blockchain
	chain := blockchain.InitBlockChain()
	// close database properly if go received exit signal
	defer func(Database *badger.DB) {
		err := Database.Close()
		if err != nil {
			// if database fails to close
			blockchain.HandleError(err)
			os.Exit(0)
		}
	}(chain.Database)

	// init cli
	cli := CommandLine{chain}
	// run cli
	cli.run()
}
