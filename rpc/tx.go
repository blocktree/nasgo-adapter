package rpc

import (
	"encoding/json"

	"github.com/assetsadapterstore/nasgo-adapter/crypto"
	"github.com/go-errors/errors"
	"gopkg.in/resty.v1"
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
		Get(tx.bk.baseAddress + "/api/transactions/get?id=" + id)
	if err != nil {
		return nil, errors.New(err)
	}
	body, err := tx.bk.ReadResponse(resp)
	response := TxResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, errors.New(err)
	}
	return response.Transaction, nil
}

func (tx *Tx) GetTransactionsByBlock(blockId string) ([]*Transaction, error) {
	resp, err := resty.
		R().
		Get(tx.bk.baseAddress + "/api/transactions?blockId=" + blockId)
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
	Result string `json:"result"`
}

func (tx *Tx) TransferNSG(txData crypto.Tx) error {
	resp, err := resty.
		R().
		SetBody(&txData).
		Post(tx.bk.baseAddress + "/peer/transactions")
	if err != nil {
		return errors.New(err)
	}
	body, err := tx.bk.ReadResponse(resp)
	if err != nil {
		return err
	}
	response := TxPublishResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return errors.New(err)
	}
	return nil
}

func (tx *Tx) TransferAsset(txData crypto.Tx) error {
	resp, err := resty.
		R().
		SetBody(&txData).
		Post(tx.bk.baseAddress + "/peer/transactions")
	if err != nil {
		return errors.New(err)
	}
	body, err := tx.bk.ReadResponse(resp)
	if err != nil {
		return err
	}
	response := TxPublishResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return errors.New(err)
	}
	return nil
}
