package accounts

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/gob"
	"jjeth/util"
	"math/big"
)

type Account struct {
	Address []byte  `json:"address"`
	Wallet  *Wallet `json:"wallet"`
	Balance big.Int `json:"balance"`
}

func GenerateAccount() *Account {
	acc := new(Account)
	acc.Balance = *big.NewInt(int64(0))
	acc.Wallet = MakeWallet()
	acc.Address = acc.Wallet.GenerateAddress()
	return acc
}

func (Acc *Account) Sign(data []byte) ([]byte, error) {
	dat := util.KeccackHash(data)
	r, s, err := ecdsa.Sign(rand.Reader, &Acc.Wallet.privKey, dat)
	if err != nil {
		return nil, err
	}
	return append(r.Bytes(), s.Bytes()...), nil
}

func VerifySignature(PubKey, data, Signature []byte) bool {
	curve := elliptic.P256()
	r := big.Int{}
	s := big.Int{}

	sigLen := len(Signature)
	r.SetBytes(Signature[:(sigLen / 2)])
	s.SetBytes(Signature[(sigLen / 2):])

	x := big.Int{}
	y := big.Int{}
	keyLen := len(PubKey)
	x.SetBytes(PubKey[:(keyLen / 2)])
	y.SetBytes(PubKey[(keyLen / 2):])
	rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
	return ecdsa.Verify(&rawPubKey, util.KeccackHash(data), &r, &s)
}

func (Acc *Account) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(Acc)
	HandleErr(err)
	return res.Bytes()
}

func DeserializeAcc(data []byte) *Account {
	acc := new(Account)
	decoder := gob.NewDecoder(bytes.NewBuffer(data))
	err := decoder.Decode(acc)
	HandleErr(err)
	return acc
}
