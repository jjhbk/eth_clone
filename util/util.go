package util

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"sort"
	"strings"

	"log"

	"github.com/google/uuid"

	"github.com/mr-tron/base58"

	geth "github.com/ethereum/go-ethereum/crypto"
)

func SortCharacters(data string) string {
	strArr := strings.Split(data, "")
	sort.Strings(strArr)
	var res string
	for _, ch := range strArr {
		res += ch
	}
	return res
}

func KeccackHash(data []byte) []byte {
	hash := geth.NewKeccakState()
	hash.Write(data)
	hs := hash.Sum(nil)
	return hs
}

func Base58Encode(input []byte) []byte {
	encode := base58.Encode(input)

	return []byte(encode)
}

func Base58Decode(input []byte) []byte {
	decode, err := base58.Decode(string(input[:]))
	if err != nil {
		log.Panic(err)
	}

	return decode
}

func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	pub := append(private.X.Bytes(), private.Y.Bytes()...)
	return *private, pub
}

func GenerateUUID() string {
	id := uuid.New()
	fmt.Println("Generated UUID:")
	fmt.Println(id.String())
	return id.String()
}
