package blockchain

import (
	"encoding/hex"
	"fmt"
	"time"
)

var JJchain *Blockchain

func init() {
	JJchain = InitChain("jj")
}

type Blockchain struct {
	Chain []*Block
}

func InitChain(beneficiary string) *Blockchain {
	genesis := new(Block)
	chain := new(Blockchain)
	genesis.Difficulty = 8
	genesis.PrevHash = []byte{}
	genesis.Height = 0
	genesis.TimeStamp = time.Now().Unix()
	chain.Chain = make([]*Block, 0)
	chain.Chain = append(chain.Chain, genesis)
	return chain
}

func (blockchain *Blockchain) Mine(beneficiary string) error {
	newBlock := MineBlock(JJchain.Chain[len(JJchain.Chain)-1], beneficiary)
	return blockchain.AddBlock(newBlock)
}

func (blockchain *Blockchain) AddBlock(block *Block) error {
	if !blockchain.ValidateBlock(block) {
		return fmt.Errorf("Error: Invalid Block")
	}
	return nil
}

func (chain *Blockchain) ValidateBlock(block *Block) bool {
	prevBlock := chain.Chain[len(chain.Chain)-1]
	h1 := hex.EncodeToString(block.PrevHash)
	h2 := hex.EncodeToString(prevBlock.Hash)
	if h1 != h2 {
		fmt.Println("prev hash error", h1, h2)
		return false
	}
	if prevBlock.TimeStamp > block.TimeStamp {
		fmt.Println("time stamp error")
		return false
	}
	if prevBlock.Height > block.Height {
		fmt.Println("block height error")

		return false
	}
	pow := NewProof(block)
	res := pow.Validate()
	if !res {
		return false
	}
	return true
}
