package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

type Block struct {
	Hash         []byte
	Transactions []*Transaction
	PRevHash     []byte
	Nonce        int
}

// HashTransactions generate hash of combiend transactions
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	// iterate over transactions
	for _, tx := range b.Transactions {
		// append to slice of bytes
		txHashes = append(txHashes, tx.ID)
	}
	// hash transactions
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

func CreateBlock(txs []*Transaction, prevHash []byte) *Block {
	block := &Block{[]byte{}, txs, prevHash, 0}
	// create new proof of work
	pow := NewProof(block)
	// run proof of work
	nonce, hash := pow.Run()

	// assign results to new block
	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func Genesis(coinbase *Transaction) *Block {
	return CreateBlock([]*Transaction{coinbase}, []byte{})
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
