package accounts_test

import (
	"jjeth/accounts"
	"jjeth/blockchain"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSignatureVerify(t *testing.T) {
	block := blockchain.Block{time.Now().Unix(), make([]*accounts.Transaction, 0), []byte{2}, []byte{4}, []byte{3}, 56, 1, 5}
	Acc := accounts.GenerateAccount()
	data := block.Serialize()
	signature, err := Acc.Sign(data)
	require.NoError(t, err)
	res := accounts.VerifySignature(Acc.Wallet.PubKey, data, signature)
	require.True(t, res)

}
