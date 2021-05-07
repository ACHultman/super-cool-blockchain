package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

// BlockChain TEMP array type
type BlockChain struct {
	blocks []*Block
}

type Block struct {
	Hash     []byte
	Data     []byte
	PRevHash []byte
}

func (b *Block) DeriveHash() {
	info := bytes.Join([][]byte{b.Data, b.PRevHash}, []byte{}) // combine data nd previous hash
	hash := sha256.Sum256(info)                                // sha256 placeholder
	b.Hash = hash[:]                                           // set hash of block ref
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash}
	block.DeriveHash()
	return block
}

func (chain *BlockChain) AddBlock(data string) {
	prevBlock := chain.blocks[len(chain.blocks)-1] // get previous block in chain
	new := CreateBlock(data, prevBlock.Hash)       // create block from data and previous block's hash
	chain.blocks = append(chain.blocks, new)       // add new block to chain
}

func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

func InitBlockChain() *BlockChain {
	return &BlockChain{[]*Block{Genesis()}}
}

func main() {
	chain := InitBlockChain()

	chain.AddBlock("First new block")
	chain.AddBlock("Second new block")
	chain.AddBlock("Third new block")

	for _, block := range chain.blocks {
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
	}
}
