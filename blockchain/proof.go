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

// Proof of work algorithm. This is basically what the miners do and get paid for
// because it is computationally expensive. This is mainly about signing Blocks
// (for security) and showing a proof that the work was done.

// Algorithm steps:
// 1. Take the data from the block
// 2. Create a counter (nonce) which starts at 0
// 3. Create a hash of the data plus the counter
// 4. Check the hash to see if it meets a set of requirements
// Requirements:
// The first few bytes must contain 0s

// Difficulty is the difficulty of the POW algorithm
const Difficulty = 12

// ProofOfWork is the main proof of work structure
type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

// NewProof creates a new ProofOfWork and returns a pointer to it
func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))

	pow := &ProofOfWork{b, target}

	return pow
}

// InitData initializes pow data
func (pow *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.Data,
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)
	return data
}

// Run runs the computation for generating the hash for block's data and check for requirements
func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0

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

	fmt.Println()

	return nonce, hash[:]
}

// Validate validates a ProofOfWork
func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int

	data := pow.InitData(pow.Block.Nonce)

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}

// ToHex converts an int64 to binary as []byte using BigEndian arrangement
func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
