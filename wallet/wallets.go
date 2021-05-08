package wallet

import (
	"blockchain/blockchain"
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
)

const walletFile = "./tmp/wallets.data"

type Wallets struct {
	Wallets map[string]*Wallet
}

func CreateWallet() (*Wallets, error) {
	// init struct
	wallets := Wallets{}
	// init wallets map
	wallets.Wallets = make(map[string]*Wallet)

	// load wallets from file
	err := wallets.LoadFile()

	return &wallets, err
}

func (ws *Wallets) AddWallet() string {
	// make wallet
	wallet := MakeWallet()
	// create address
	address := fmt.Sprintf("%s", wallet.Address())

	// add wallet to wallets with address as key
	ws.Wallets[address] = wallet

	return address
}

// GetAllAddresses returns all wallet addresses in wallets
func (ws *Wallets) GetAllAddresses() []string {
	var addresses []string

	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

func (ws Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

func (ws *Wallets) LoadFile() error {
	// check if file exists
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}

	var wallets Wallets

	// read bytes from file
	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		return err
	}

	// register elliptic curve
	gob.Register(elliptic.P256())
	// init decoder for file bytes
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	// decode bytes
	err = decoder.Decode(&wallets)

	// put docded wallets into struct
	ws.Wallets = wallets.Wallets

	return nil
}

func (ws *Wallets) SaveFile() {
	var content bytes.Buffer

	// register elliptic curve
	gob.Register(elliptic.P256())

	// create encoder
	encoder := gob.NewEncoder(&content)
	// encode data
	err := encoder.Encode(ws)
	blockchain.HandleError(err)

	// write encoded data to file
	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644) // rw permisisons
	blockchain.HandleError(err)
}
