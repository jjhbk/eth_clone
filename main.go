package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"jjeth/blockchain"
	"jjeth/network"
	"sync"
)

var Mutex = &sync.Mutex{}
var Blockchain *blockchain.Blockchain

func main() {

	//chain := blockchain.InitChain("jj", map[string]int{})
	//

	target := flag.String("d", "", "target peer to dial")
	listenF := flag.Int("l", 0, "wait for incoming connections")
	blocks := flag.Int("b", 1, "enter number of blocks to be mined")
	flag.Parse()
	if *blocks > 0 {
		for i := 1; i < *blocks; i++ {
			block := blockchain.MineBlock(blockchain.JJchain.Chain[len(blockchain.JJchain.Chain)-1], "jj")

			blockchain.JJchain.Chain = append(blockchain.JJchain.Chain, block)
			fmt.Printf("%v", block)
		}
	}
	hash := blockchain.JJchain.Chain[1].Hash
	fmt.Println(len(hash), hex.EncodeToString(hash), string(hash))
	network.StartServer(*target, *listenF)

}
