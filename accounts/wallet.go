package accounts

import (
	"bytes"
	"crypto/ecdsa"
	"jjeth/util"
	"log"

	"golang.org/x/crypto/ripemd160"
)

const (
	checksumLength = 4
	version        = byte(0x00)
)

type Wallet struct {
	privKey ecdsa.PrivateKey `json:"privkey"`
	PubKey  []byte           `json:"pubkey"`
}

func (w Wallet) GenerateAddress() []byte {
	pubHash := PublicKeyHash(w.PubKey)
	versionedHash := append([]byte{version}, pubHash...)
	checksum := Checksum(versionedHash)
	fullHash := append(versionedHash, checksum...)
	address := util.Base58Encode(fullHash)
	return address
}

func MakeWallet() *Wallet {
	private, public := util.NewKeyPair()
	wallet := Wallet{private, public}
	return &wallet
}

func PublicKeyHash(pubKey []byte) []byte {
	pubHash := util.KeccackHash(pubKey)
	hasher := ripemd160.New()
	_, err := hasher.Write(pubHash[:])
	if err != nil {
		log.Panic(err)
	}
	publicRipMD := hasher.Sum(nil)
	return publicRipMD
}

func Checksum(payload []byte) []byte {
	firstHash := util.KeccackHash(payload)
	secondHash := util.KeccackHash(firstHash[:])
	return secondHash[:checksumLength]
}

func ValidateAddress(address string) bool {
	pubKeyHash := util.Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-checksumLength:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-checksumLength]
	targetCheckSum := Checksum(append([]byte{version}, pubKeyHash...))
	return bytes.Compare(actualChecksum, targetCheckSum) == 0
}
