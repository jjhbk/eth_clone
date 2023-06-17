package accounts_test

import (
	"fmt"
	"jjeth/accounts"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTrei(t *testing.T) {
	trie := accounts.NewTrie()
	trie.Put("foo", []byte("jjhbk"))
	val := trie.Get("foo")

	fmt.Printf("%v", trie)

	require.Equal(t, val, []byte("jjhbk"))
}
