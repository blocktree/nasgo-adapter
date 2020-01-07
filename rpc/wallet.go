package rpc

import (
	"encoding/json"

	"github.com/go-errors/errors"
	"gopkg.in/resty.v1"
)

type Wallet struct {
	bk *BaseClient
}

type BalanceResponse struct {
	Balance uint64 `json:"balance"`
}

type AssetsBalanceResponse struct {
	Balance *AssetsBalance `json:"balance"`
}

type AssetsBalance struct {
	Currency  string `json:"currency"`
	Balance   string `json:"balance"`
	Precision uint8  `json:"precision"`
}

func newWalletClient(bk *BaseClient) *Wallet {
	return &Wallet{
		bk: bk,
	}
}

func (w *Wallet) GetBalance(address string) (uint64, error) {
	resp, err := resty.
		R().
		Get(w.bk.baseAddress + "/api/accounts/getBalance?address=" + address)
	if err != nil {
		return 0, err
	}
	body, err := w.bk.ReadResponse(resp)
	if err != nil {
		return 0, err
	}
	balanceResponse := BalanceResponse{}
	if err := json.Unmarshal(body, &balanceResponse); err != nil {
		return 0, errors.New(err)
	}
	return balanceResponse.Balance, nil
}

func (w *Wallet) GetAssetsBalance(address, currency string) (*AssetsBalance, error) {
	resp, err := resty.
		R().
		Get(w.bk.baseAddress + "/api/uia/balances/" + address + "/" + currency)
	if err != nil {
		return nil, err
	}
	body, err := w.bk.ReadResponse(resp)
	if err != nil {
		return nil, err
	}
	balanceResponse := AssetsBalanceResponse{}
	if err := json.Unmarshal(body, &balanceResponse); err != nil {
		return nil, errors.New(err)
	}
	return balanceResponse.Balance, nil
}
