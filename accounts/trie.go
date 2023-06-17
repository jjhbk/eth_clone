package accounts

import (
	"encoding/json"
	"fmt"
	"jjeth/util"
	"log"
	"strings"
)

type Trie struct {
	Head     *Node  `json:"head"`
	RootHash []byte `json:"roothash"`
}

type Node struct {
	Value    []byte           `json:"value"`
	ChildMap map[string]*Node `json:"childmap"`
}

func NewTrie() *Trie {
	trie := new(Trie)
	trie.Head = NewNode()
	trie.RootHash = trie.GenerateRootHash()
	return trie
}

func NewNode() *Node {
	node := new(Node)
	node.ChildMap = map[string]*Node{}
	return node
}
func (trie *Trie) Get(key string) []byte {
	if trie.Head == nil {
		log.Panic("nil trie error")
	}
	node := trie.Head
	charr := strings.Split(key, "")
	for _, ch := range charr {
		if node.ChildMap[ch] != nil {
			node = node.ChildMap[ch]
		} else {
			return nil
		}
	}
	return node.Value
}

func (trie *Trie) Put(key string, value []byte) {
	if trie == nil {
		trie = NewTrie()
	}
	node := trie.Head
	charr := strings.Split(key, "")
	for i, ch := range charr {
		fmt.Println(i)
		val, ok := node.ChildMap[ch]
		if !ok {
			node.ChildMap[ch] = NewNode()
			val = node.ChildMap[ch]
		}
		node = val
	}
	node.Value = value
	trie.RootHash = trie.GenerateRootHash()
}

func (trie *Trie) GenerateRootHash() []byte {
	dat, err := json.Marshal(trie)
	if err != nil {
		log.Panic("marshalling error")
	}
	return util.KeccackHash(dat)
}
