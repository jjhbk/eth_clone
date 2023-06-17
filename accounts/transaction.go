package accounts

import (
	"bytes"
	"encoding/gob"
	"jjeth/util"
	"log"
	"math/big"
)

const (
	CREATE_ACCOUNT = "CREATEACCOUNT"
	TRANSACT       = "TRANSACT"
)

var TxQueue *TransactionQueue

func init() {
	TxQueue = new(TransactionQueue)
	TxQueue.txMap = make(map[string]*Transaction)
}

type Transaction struct {
	Id        string  `json:"id"`
	From      []byte  `json:"from"`
	Type      string  `json:"type"`
	To        []byte  `json:"to"`
	Value     big.Int `json:"value"`
	Data      []byte  `json:"data"`
	Signature []byte  `json:"signature"`
}

type TransactionQueue struct {
	txMap map[string]*Transaction
}

func (txqueue *TransactionQueue) AddTx(tx *Transaction) {
	txqueue.txMap[tx.Id] = tx
}

func (txqueue *TransactionQueue) GetTransactionSeries() map[string]*Transaction {
	return txqueue.txMap
}

func CreateTransaction(Acc *Account, to []byte, value big.Int, typ string) *Transaction {
	transaction := Transaction{Id: util.GenerateUUID(), From: Acc.Wallet.PubKey, To: to, Value: value, Type: typ, Data: []byte("TRANSACTION DATA"), Signature: []byte{}}
	transactionData := transaction.Serialize()
	var err error
	transaction.Signature, err = Acc.Sign(transactionData)
	HandleErr(err)
	TxQueue.AddTx(&transaction)
	return &transaction
}

func (tx *Transaction) ValidateTransaction() bool {
	signature := tx.Signature
	txCopy := *tx
	txCopy.Signature = []byte{}
	return VerifySignature(tx.From, txCopy.Serialize(), signature)
}
func (b *Transaction) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(b)
	HandleErr(err)
	return res.Bytes()
}

func ClearBlockTransactions(txList []*Transaction) {
	for _, tx := range txList {
		delete(TxQueue.txMap, tx.Id)
	}
}

func DeserializeTx(data []byte) *Transaction {
	var Transaction Transaction
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&Transaction)
	HandleErr(err)
	return &Transaction
}

func (tx *Transaction) RunTransaction() {
	switch tx.Type {
	case TRANSACT:
		tx.runStandardTransaction()
	case CREATE_ACCOUNT:
		tx.runCreateAccountTransaction()
	default:
		break
	}
}

func (tx *Transaction) runStandardTransaction() {
	fromAccount := JJETH_STATE.stateTrie.Get(string(tx.From))
	toAccount := JJETH_STATE.stateTrie.Get(string(tx.To))
	if fromAccount == nil {
		log.Panic("error: account does not exist")
	}
	if toAccount == nil {

	}
	fromAcc := DeserializeAcc(fromAccount)
	toAcc := DeserializeAcc(toAccount)
	fromAcc.Balance.Sub(&fromAcc.Balance, &tx.Value)
	toAcc.Balance.Add(&toAcc.Balance, &tx.Value)
	fromAccount = fromAcc.Serialize()
	toAccount = toAcc.Serialize()
	JJETH_STATE.putAccount(string(tx.From), fromAccount)
	JJETH_STATE.putAccount(string(tx.To), toAccount)
}

func (tx *Transaction) runCreateAccountTransaction() {
	accData := tx.Data
	acc := GenerateAccount()
	acc.Address = accData
	JJETH_STATE.putAccount(string(acc.Address), acc.Serialize())

}

func HandleErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
