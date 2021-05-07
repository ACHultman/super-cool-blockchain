package main

import (
	"blockchain/blockchain"
	"flag"
	"fmt"
	"github.com/dgraph-io/badger"
	"log"
	"os"
	"runtime"
	"strconv"
)

type CommandLine struct {
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" getBalance -address ADDRESS - get the balance for address")
	fmt.Println(" createblockchain -address ADRESS - creates a blockchain")
	fmt.Println(" send -from FROM -to TO -amount AMOUNT - Send amount between addresses")
	fmt.Println(" printchain - Prints the blocks in the chain")
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

func (cli *CommandLine) printChain() {
	chain := blockchain.ContinueBlockChain()
	defer func(Database *badger.DB) {
		err := Database.Close()
		if err != nil {
			os.Exit(1)
		}
	}(chain.Database)
	iter := chain.Iterator()

	for {
		block := iter.Next()
		fmt.Printf("Prev Hash: %s\n", block.PRevHash)
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

func (cli *CommandLine) createBlockChain(address string) {
	chain := blockchain.InitBlockChain(address)
	// close database properly if go received exit signal
	err := chain.Database.Close()
	if err != nil {
		return
	}
	fmt.Println("Finished!")
}

func (cli *CommandLine) getBalance(address string) {
	chain := blockchain.ContinueBlockChain()
	// close database properly if go received exit signal
	defer chain.Database.Close()

	balance := 0
	UTXOs := chain.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d\n", address, balance)
}

func (cli *CommandLine) send(from, to string, amount int) {
	chain := blockchain.ContinueBlockChain()
	defer func(Database *badger.DB) {
		err := Database.Close()
		if err != nil {
			blockchain.HandleError(err)
			return
		}
	}(chain.Database)

	// create new transaction
	tx := blockchain.NewTransaction(from, to, amount, chain)
	// add block containing transaction
	chain.AddBlock([]*blockchain.Transaction{tx})
	fmt.Println("Success!")
}

func (cli *CommandLine) run() {
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			runtime.Goexit()
		}
		cli.createBlockChain(*createBlockchainAddress)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}

		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
}

func main() {
	// help to properly close database
	defer os.Exit(0)
	// init cli
	cli := CommandLine{}
	// run cli
	cli.run()
}
