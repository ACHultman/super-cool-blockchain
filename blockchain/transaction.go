package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

type Transaction struct {
	ID        []byte
	Inputs    []TxInput
	Outputess []TxOutput
}

type TxOutput struct {
	Value  int
	PubKey string
}

type TxInput struct {
	ID  []byte
	Out int
	Sig string
}

func (tx *Transaction) SetID() {
	// init bytes buffer
	var encoded bytes.Buffer
	var hash [32]byte

	// init encoder
	encoder := gob.NewEncoder(&encoded)
	// encode transaction
	err := encoder.Encode(tx)
	HandleError(err)

	// hash encoded bytes with sha256
	hash = sha256.Sum256(encoded.Bytes())
	// set transaction ID as hash
	tx.ID = hash[:]
}

func CoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Coins to %s", to)
	}

	// init input, output
	// empty in
	txin := TxInput{[]byte{}, -1, data}
	// constant reward of 100
	txout := TxOutput{100, to}

	// init transaction struct
	tx := Transaction{nil, []TxInput{txin}, []TxOutput{txout}}
	// create hash ID for transaction
	tx.SetID()

	return &tx
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

func (in *TxInput) CanUnlock(data string) bool {
	return in.Sig == data
}

func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.PubKey == data
}
