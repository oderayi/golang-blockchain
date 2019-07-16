package blockchain

import (
	"fmt"

	"github.com/dgraph-io/badger"
)

const (
	dbPath = "./tmp/blocks"
)

// BlockChain is the blockchain data structure
type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

//BlockChainIterator is an iterator for BlockChain
type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

// InitBlockChain initializes the blockchain
func InitBlockChain() *BlockChain {
	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)

	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found")
			genesis := Genesis()
			fmt.Println("Genesis proved")
			err = txn.Set(genesis.Hash, genesis.Serialize())

			err = txn.Set([]byte("lh"), genesis.Hash)

			lastHash = genesis.Hash

			return err
		}

		item, err := txn.Get([]byte("lh"))
		Handle(err)
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		return err
	})

	Handle(err)

	blockchain := BlockChain{lastHash, db}

	return &blockchain
}

// AddBlock adds a new Block to the blockchain
func (chain *BlockChain) AddBlock(data string) {
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})

		return err
	})

	Handle(err)

	newBlock := CreateBlock(data, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash

		return err
	})
	Handle(err)
}

// Iterator returns a BlockChainIterator for iterating over BlockChain
func (chain *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{chain.LastHash, chain.Database}

	return iter
}

// Next returns the "next" Block in the BlockChain (backward iteration)
func (iter *BlockChainIterator) Next() *Block {
	var block *Block

	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)

		err = item.Value(func(val []byte) error {
			encodedBlock := val
			block = Deserialize(encodedBlock)
			return nil
		})

		return err
	})
	Handle(err)

	iter.CurrentHash = block.PrevHash

	return block
}
