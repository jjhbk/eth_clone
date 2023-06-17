package blockchain

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"jjeth/util"
	"log"
	"math"
	"math/big"
)

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProof(Block *Block) *ProofOfWork {
	num := big.NewInt(1)
	target := num.Lsh(num, (256 - uint(Block.Difficulty)))
	return &ProofOfWork{Target: target, Block: Block}
}

func (pow *ProofOfWork) InitData(nonce int) []byte {
	//data := bytes.Join([][]byte{pow.Block.PrevHash, ToHex(int64(nonce)), ToHex(int64(pow.Block.Difficulty))}, []byte{})
	pow.Block.Nonce = nonce
	data := pow.Block.Serialize()
	return data
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash, finalHash []byte
	nonce := 0
	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		hash = util.KeccackHash(data)
		fmt.Printf("\r%x", hash)
		finalHash = hash[:]
		intHash.SetBytes(hash[:])

		if intHash.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}
	}
	return nonce, finalHash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int
	blockhash := pow.Block.Hash
	pow.Block.Hash = []byte{}
	data := pow.InitData(pow.Block.Nonce)
	hash := util.KeccackHash(data)
	intHash.SetBytes(hash[:])
	pow.Block.Hash = blockhash
	return intHash.Cmp(pow.Target) == -1
}

func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)

	}

	return buff.Bytes()
}
