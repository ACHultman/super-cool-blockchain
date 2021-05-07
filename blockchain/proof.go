package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

// Difficulty TODO increase Difficulty algorithmicaly so it stays hard with increasing miners and computational power
const Difficulty = 18

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

// NewProof Pair block with target to create proof of work struct
func NewProof(b *Block) *ProofOfWork {
	// Init target
	target := big.NewInt(1)
	// left shift decided by Difficulty
	target.Lsh(target, uint(256-Difficulty))

	// Pair block with target
	pow := &ProofOfWork{b, target}

	return pow
}

// InitData hash data using previous hash, nonce, difficulty
func (pow *ProofOfWork) InitData(nonce int) []byte {
	// combine Previous Hash and Data into byte struct
	data := bytes.Join(
		[][]byte{
			pow.Block.PRevHash,
			pow.Block.HashTransactions(),
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)
	return data
}

// run proof of work, compare target with hashes of data with nonce
func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0

	// effectively infinite loop
	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		intHash.SetBytes(hash[:])

		if intHash.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}
	}

	fmt.Print()

	return nonce, hash[:]
}

// Validate compare nonce hashed data with the target to verify
func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int

	data := pow.InitData(pow.Block.Nonce)

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}

// ToHex convert int64 to slice of bytes
func ToHex(num int64) []byte {
	// create new buffer
	buff := new(bytes.Buffer)
	// write num to buffer
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		// if error
		log.Panic(err)
	}

	return buff.Bytes()
}
