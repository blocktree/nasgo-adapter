package rpc

import (
	"encoding/json"
	"github.com/blocktree/openwallet/log"
	"github.com/imroc/req"

	"github.com/go-errors/errors"
	"gopkg.in/resty.v1"
)

const (
	TxType_NSG            = 0  //	NSG Transactions
	TxType_SetSecureCode  = 1  //Set secure code
	TxType_Delegate       = 2  //Register as a delegator
	TxType_Vote           = 3  //Submit a vote
	TxType_MultiSig       = 4  //multisignature
	TxType_PublishDAPP    = 5  //Publish a dapp in mainnet
	TxType_DeopsitDAPP    = 6  //deposit to a Dapp
	TxType_WithdrawalDAPP = 7  //withdrawal from a dapp
	TxType_RegPublisher   = 9  //register as a asset publisher
	TxType_RegAsset       = 10 //register an asset in mainnet
	TxType_IssueAsset     = 13 //issue an asset
	TxType_Asset          = 14 //asset transactions
)

type TxResponse struct {
	Success     bool         `json:"success"`
	Transaction *Transaction `json:"transaction"`
}

type TxsResponse struct {
	Success      bool           `json:"success"`
	Transactions []*Transaction `json:"transactions"`
}

type Transaction struct {
	ID            string   `json:"id"`
	Height        string   `json:"height"`
	BlockID       string   `json:"blockId"` // block id
	Type          uint32   `json:"type"`
	Timestamp     int64    `json:"timestamp"` // A timestamp recording when this block was created (Will overflow in 2106[2])
	SenderID      string   `json:"senderId"`
	RecipientId   string   `json:"recipientId"`
	Amount        uint64   `json:"amount"`
	Fee           uint64   `json:"fee"`
	Signature     string   `json:"signature"`
	Signatures    []string `json:"signatures"`
	SignSignature string   `json:"signSignature"`
	Confirmations string   `json:"confirmations"`
	Message       string   `json:"message"`
	Asset         *Asset   `json:"asset"`
}

type Asset struct {
	TransactionId string `json:"transactionId"`
	Currency      string `json:"currency"`
	Amount        string `json:"amount"`
	Precision     uint8  `json:"precision"`
}

type Tx struct {
	bk *BaseClient
}

func newTxClient(bk *BaseClient) *Tx {
	return &Tx{
		bk: bk,
	}
}

func (tx *Tx) GetTransaction(id string) (*Transaction, error) {
	resp, err := resty.
		R().
		Get(tx.bk.baseAddress + "/api/uia/transactions/get?id=" + id)
	if err != nil {
		return nil, errors.New(err)
	}
	body, err := tx.bk.ReadResponse(resp)
	response := TxsResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, errors.New(err)
	}
	return response.Transactions[0], nil
}

func (tx *Tx) GetTransactionsByBlock(blockId string) ([]*Transaction, error) {
	resp, err := resty.
		R().
		Get(tx.bk.baseAddress + "/api/uia/transactions?blockId=" + blockId)
	if err != nil {
		return nil, errors.New(err)
	}
	body, err := tx.bk.ReadResponse(resp)
	response := TxsResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, errors.New(err)
	}
	return response.Transactions, nil
}

type TxPublishResponse struct {
	Result bool `json:"success"`
}

func (tx *Tx) Broadcast(txData interface{}) error {
	b, err := json.Marshal(txData)
	if err != nil {
		return err
	}
	log.Debugf("Broadcast tx: %s", string(b))
	resp, err := resty.
		R().
		SetBody(b).
		SetHeader("version", "''").
		SetHeader("magic", "594fe0f3").
		Post(tx.bk.baseAddress + "/peer/transactions")
	if err != nil {
		return errors.New(err)
	}
	body, err := tx.bk.ReadResponse(resp)
	if err != nil {
		return err
	}
	response := TxPublishResponse{}
	if err := json.Unmarshal(body, &response); err != nil || !response.Result {
		return errors.New(err)
	}
	return nil
}


func (tx *Tx) BroadcastTx(txData interface{}) error {

	authHeader := req.Header{
		"version": "''",
		"magic": "594fe0f3",
		"Content-Type":  "application/json",
	}

	r, err := req.Post(tx.bk.baseAddress + "/peer/transactions", req.BodyJSON(txData), authHeader)
	if err != nil {
		return errors.New(err)
	}
	log.Std.Info("%+v", r)
	log.Debugf("response: %s", r.String())
	return nil
}
