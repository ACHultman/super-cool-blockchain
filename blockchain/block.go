package blockchain

// BlockChain TEMP array type
type BlockChain struct {
	Blocks []*Block
}

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

func (chain *BlockChain) AddBlock(data string) {
	prevBlock := chain.Blocks[len(chain.Blocks)-1] // get previous block in chain
	new := CreateBlock(data, prevBlock.Hash)       // create block from data and previous block's hash
	chain.Blocks = append(chain.Blocks, new)       // add new block to chain
}

func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

func InitBlockChain() *BlockChain {
	return &BlockChain{[]*Block{Genesis()}}
}
