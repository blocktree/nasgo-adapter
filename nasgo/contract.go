/*
 * Copyright 2018 The OpenWallet Authors
 * This file is part of the OpenWallet library.
 *
 * The OpenWallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The OpenWallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */

package nasgo

import (
	"github.com/blocktree/openwallet/openwallet"
	"github.com/shopspring/decimal"
)

type ContractDecoder struct {
	openwallet.SmartContractDecoderBase
	wm *WalletManager
}

//NewContractDecoder 智能合约解析器
func NewContractDecoder(wm *WalletManager) *ContractDecoder {
	decoder := ContractDecoder{}
	decoder.wm = wm
	return &decoder
}

// GetTokenBalanceByAddress return the balance by address, queried by rpc
func (decoder *ContractDecoder) GetTokenBalanceByAddress(contract openwallet.SmartContract, address ...string) ([]*openwallet.TokenBalance, error) {

	tokenBalanceList := make([]*openwallet.TokenBalance, 0)

	for _, addr := range address {

		balance, err := decoder.wm.WalletClient.Wallet.GetAssetsBalance(addr, contract.Address)
		if err != nil {
			decoder.wm.Log.Errorf("get account[%v] token balance failed, err: %v", addr, err)
		}

		if balance == nil {
			return nil, err
		}

		value, _ := decimal.NewFromString(balance.Balance)
		value = value.Shift(-int32(contract.Decimals))
		// value = value.Shift(-int32(balance.Precision))

		tokenBalance := &openwallet.TokenBalance{
			Contract: &contract,
			Balance: &openwallet.Balance{
				Address:          addr,
				Symbol:           contract.Symbol,
				Balance:          value.String(),
				ConfirmBalance:   value.String(),
				UnconfirmBalance: "0",
			},
		}

		tokenBalanceList = append(tokenBalanceList, tokenBalance)
	}

	return tokenBalanceList, nil

}
