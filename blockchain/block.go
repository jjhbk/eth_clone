package blockchain

import (
	"bytes"
	"encoding/gob"
	"jjeth/accounts"
	"log"
	"time"
)

const MINE_RATE = 13

type Block struct {
	TimeStamp    int64                   `json:"block"`
	Transactions []*accounts.Transaction `json:"transactions"`
	Hash         []byte                  `json:"hash"`
	PrevHash     []byte                  `json:"prevhash"`
	StateRoot    []byte                  `json:"stateroot"`
	Nonce        int                     `json:"nonce"`
	Height       int                     `json:"height"`
	Difficulty   int                     `json:"difficulty"`
}

func adjustDifficulty(lastBlock *Block, time int64) int {
	if (time-lastBlock.TimeStamp) > MINE_RATE && lastBlock.Difficulty > 1 {
		return lastBlock.Difficulty - 1
	}
	return lastBlock.Difficulty + 1
}

func MineBlock(lastBlock *Block, beneficiary string) *Block {
	timestamp := time.Now().Unix()
	difficulty := adjustDifficulty(lastBlock, timestamp)
	block := &Block{timestamp, make([]*accounts.Transaction, 0), []byte{}, lastBlock.Hash, []byte{}, 0, lastBlock.Height + 1, difficulty}
	for _, tx := range accounts.TxQueue.GetTransactionSeries() {
		block.Transactions = append(block.Transactions, tx)
	}
	block.StateRoot = accounts.JJETH_STATE.GetStateRoot()
	block.RunBlock()
	pow := NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce
	if !pow.Validate() {
		log.Panic("error validating proof of work")
	}
	accounts.ClearBlockTransactions(block.Transactions)
	return block
}

func (b *Block) RunBlock() {
	for _, tx := range b.Transactions {
		tx.RunTransaction()
	}
}
func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(b)
	HandleErr(err)
	return res.Bytes()
}

func Deserialize(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	HandleErr(err)
	return &block
}
func HandleErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
