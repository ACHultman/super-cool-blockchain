package blockchain

import (
	"encoding/hex"
	"fmt"
	"github.com/dgraph-io/badger"
	"os"
	"runtime"
)

const (
	dbPath      = "./tmp/blocks"
	dbFile      = "./tmp/blocks/MANIFEST"
	genesisData = "First Transaction"
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

// DBExists checks if DB MANIFEST file exists
func DBExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

func ContinueBlockChain() *BlockChain {
	if DBExists() == false {
		fmt.Println("No existing blockchaqin found, create one!")
		runtime.Goexit()
	}

	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	HandleError(err)
	err = db.Update(func(txn *badger.Txn) error {
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
	})
	HandleError(err)

	chain := BlockChain{lastHash, db}

	return &chain
}

func InitBlockChain(address string) *BlockChain {
	var lastHash []byte

	if DBExists() {
		fmt.Println("Blockchain alreadyt exists")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	HandleError(err)

	err = db.Update(func(txn *badger.Txn) error {
		cbtx := CoinbaseTx(address, genesisData)
		genesis := Genesis(cbtx)
		fmt.Println("Genesis created")
		err = txn.Set(genesis.Hash, genesis.Serialize())
		HandleError(err)
		err = txn.Set([]byte("lh"), genesis.Hash)

		lastHash = genesis.Hash

		return err
	})
	HandleError(err)

	blockchain := BlockChain{lastHash, db}
	return &blockchain
}

func (chain *BlockChain) AddBlock(transactions []*Transaction) {
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

	newBlock := CreateBlock(transactions, lastHash)

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

func (chain *BlockChain) FindUnspentTransactions(address string) []Transaction {
	var unspentTxs []Transaction

	// init key string, value []int map
	spentTXOs := make(map[string][]int)

	// init iterator
	iter := chain.Iterator()

	// loop until Genesis block
	for {
		// get next block
		block := iter.Next()

		// loop over block's transactions
		for _, tx := range block.Transactions {
			// get transaction ID as hex string
			txID := hex.EncodeToString(tx.ID)

			// iterate over outputs inside transaction
		Outputs:
			for outIDx, out := range tx.Outputs {
				// if output is not inside map
				if spentTXOs[txID] != nil {
					// iterate over map
					for _, spentOut := range spentTXOs[txID] {
						// check if spent out is this out
						if spentOut == outIDx {
							continue Outputs
						}
					}
				}
				// check if out belongs to address
				if out.CanBeUnlocked(address) {
					// add to unspent transactions
					unspentTxs = append(unspentTxs, *tx)
				}
			}
			// check if transaction is a coinbase transaction
			if tx.IsCoinbase() == false {
				// not a coinbase transaction
				// iterate over inputs
				for _, in := range tx.Inputs {
					// check if input belongs to address
					if in.CanUnlock(address) {
						// can unlock with address
						// convert input ID ot hex string
						inTxID := hex.EncodeToString(in.ID)
						// append to map
						spentTXOs[inTxID] = append(spentTXOs[inTxID])
					}
				}
			}
		}
		if len(block.PRevHash) == 0 {
			break
		}
	}
	return unspentTxs
}

func (chain *BlockChain) FindUTXO(address string) []TxOutput {
	var UTXOs []TxOutput
	unspentTransactions := chain.FindUnspentTransactions(address)

	// iterate over unspent transactions
	for _, tx := range unspentTransactions {
		// iterate over transaction outputs
		for _, out := range tx.Outputs {
			// if output belongs to address
			if out.CanBeUnlocked(address) {
				// add to unspent transaction outputs
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

func (chain *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxs := chain.FindUnspentTransactions(address)
	accumulated := 0

	// iterate over unspent transactions
Work:
	for _, tx := range unspentTxs {
		// convert transcation ID to hex string
		txID := hex.EncodeToString(tx.ID)
		// iterate over unspent transaction outputes
		for outIdx, out := range tx.Outputs {
			// check if output belongs to address and if user has enough tokens
			if out.CanBeUnlocked(address) && accumulated < amount {
				// increment accumulated value by output value
				accumulated += out.Value
				// add unspect output to map
				unspentOuts[txID] = append(unspentOuts[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOuts
}
