package wallet

import (
	"blockchain/blockchain"
	"github.com/mr-tron/base58"
)

// base58 == base64 without chars: 0, O, 1, I, +, /

func Base58Encode(input []byte) []byte {
	encode := base58.Encode(input)
	return []byte(encode)
}

func Base58Decode(input []byte) []byte {
	decode, err := base58.Decode(string(input[:]))
	blockchain.HandleError(err)
	return decode
}
