package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Block struct {
	Hash     []byte
	Data     []byte
	PRevHash []byte
	Nonce    int
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash, 0}
	// create new proof of work
	pow := NewProof(block)
	// run proof of work
	nonce, hash := pow.Run()

	// assign results to new block
	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

func (b *Block) Serialize() []byte {
	// init bytes buffer
	var res bytes.Buffer

	// init encoder
	encoder := gob.NewEncoder(&res)

	// encode block
	err := encoder.Encode(b)

	HandleError(err)

	// return bytes portion of block
	return res.Bytes()
}

func Deserialize(data []byte) *Block {
	// init block
	var block Block

	// init decoder
	decoder := gob.NewDecoder(bytes.NewReader(data))

	// decode into block
	err := decoder.Decode(&block)

	HandleError(err)

	return &block
}

func HandleError(err error) {
	if err != nil {
		log.Panic(err)
	}
}
