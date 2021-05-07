package main

import (
	"blockchain/blockchain"
	"fmt"
	"strconv"
)

func main() {
	chain := blockchain.InitBlockChain()

	chain.AddBlock("First new block")
	chain.AddBlock("Second new block")
	chain.AddBlock("Third new block")

	for _, block := range chain.Blocks {
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := blockchain.NewProof(block)
		fmt.Printf("POW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}
