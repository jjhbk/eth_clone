package accounts_test

import (
	"fmt"
	"jjeth/accounts"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateTransaction(t *testing.T) {
	acc := accounts.GenerateAccount()
	tx := accounts.CreateTransaction(acc, []byte("jj"), *big.NewInt(int64(100)), accounts.TRANSACT)
	fmt.Println(tx)
	require.True(t, tx.ValidateTransaction())
}
