package blockchain

import (
	"fmt"
	"github.com/dgraph-io/badger"
)

const (
	dbPath = "./tmp/blocks"
)

// BlockChain TEMP array type
type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func InitBlockChain() *BlockChain {
	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)

	HandleError(err)

	err = db.Update(func(txn *badger.Txn) error {
		// check for last hash (db/blockhain)
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			// If no last hash found
			fmt.Println("No existing blockchain found! Creating one...")
			genesis := Genesis()
			fmt.Println("Genesis proved")
			err = txn.Set(genesis.Hash, genesis.Serialize())
			HandleError(err)
			err = txn.Set([]byte("lh"), genesis.Hash)

			lastHash = genesis.Hash

			return err
		} else {
			// If last hash exists (db, blockchain exists)
			// Get last hash item
			item, err := txn.Get([]byte("lh"))
			HandleError(err)
			// Get last hash value
			err = item.Value(func(val []byte) error {
				// This func with val would only be called if item.Value encounters no error.
				// Copying or parsing val is valid.
				lastHash = append([]byte{}, val...)
				fmt.Printf("Last Hash is: %x\n", lastHash)
				return nil
			})
			return err
		}
	})

	HandleError(err)

	blockchain := BlockChain{lastHash, db}
	return &blockchain
}

func (chain *BlockChain) AddBlock(data string) {
	var lastHash []byte
	// get previous block in chain

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		HandleError(err)
		err = item.Value(func(val []byte) error {
			// This func with val would only be called if item.Value encounters no error.
			// Copying or parsing val is valid.
			lastHash = append([]byte{}, val...)
			fmt.Printf("Last Hash is: %x\n", lastHash)
			return nil
		})
		return err
	})
	HandleError(err)

	newBlock := CreateBlock(data, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		// insert block into db: {key=hash, value=serialized data}
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		HandleError(err)
		// update last hash key value as current block hash
		err = txn.Set([]byte("lh"), newBlock.Hash)

		// update blockchain struct's last hash field with current block's hash
		chain.LastHash = newBlock.Hash

		return err
	})
	HandleError(err)
}

// Iterator construct new iterator from blockchain's last hash and database
func (chain *BlockChain) Iterator() *BlockChainIterator {
	// will be iterating from last hash to genesis, so use blockchain's last hash to initialzie iterator
	iter := &BlockChainIterator{chain.LastHash, chain.Database}
	return iter
}

func (iter *BlockChainIterator) Next() *Block {
	var block *Block

	err := iter.Database.View(func(txn *badger.Txn) error {
		// get serialized current hash
		item, err := txn.Get(iter.CurrentHash)
		HandleError(err)
		err = item.Value(func(val []byte) error {
			// This func with val would only be called if item.Value encounters no error.
			// Copy and deserialize encoded val
			block = Deserialize(append([]byte{}, val...))
			fmt.Printf("Last Hash is: %x\n", block.Hash)
			return nil
		})
		return err
	})
	HandleError(err)
	// update iterator's current hash
	iter.CurrentHash = block.PRevHash
	return block
}
