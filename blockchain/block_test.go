package blockchain_test

import (
	"fmt"
	"jjeth/blockchain"
	"testing"
)

func TestMineBlock(t *testing.T) {
	block := new(blockchain.Block)
	block.Difficulty = 1
	block.Hash = []byte{1}
	newBlock := blockchain.MineBlock(block, "jj")
	for i := 0; i < 5; i++ {
		block := blockchain.MineBlock(newBlock, "jj")
		newBlock = block
		fmt.Println(newBlock)
	}

}
