/*
 * Copyright 2020 The openwallet Authors
 * This file is part of the openwallet library.
 *
 * The openwallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The openwallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */

package nasgo

import (
	"github.com/blocktree/openwallet/v2/log"
	"github.com/blocktree/openwallet/v2/openwallet"
	"path/filepath"
	"testing"

	"github.com/astaxie/beego/config"
)

var (
	tw *WalletManager
)

func init() {

	tw = testNewWalletManager()
}

func testNewWalletManager() *WalletManager {
	wm := NewWalletManager()

	//读取配置
	absFile := filepath.Join("conf", "NSG.ini")
	//log.Debug("absFile:", absFile)
	c, err := config.NewConfig("ini", absFile)
	if err != nil {
		return nil
	}
	wm.LoadAssetsConfig(c)
	return wm
}

func TestContractDecoder_GetTokenBalanceByAddress(t *testing.T) {
	contract := openwallet.SmartContract{
		ContractID: "",
		Symbol:     "NSG",
		Address:    "IMM.IMM",
		Token:      "IMMT",
		Protocol:   "",
		Name:       "IMMT",
		Decimals:   0,
	}
	addr := "NP2YbwgZHCY9tEnUVcfUmQmzUCun2wJ17F"
	tokens, err := tw.ContractDecoder.GetTokenBalanceByAddress(contract, addr)
	if err != nil {
		t.Errorf("GetTokenBalanceByAddress failed, err: %v", err)
		return
	}
	for _, t := range tokens {
		log.Infof("token: %+v", t.Balance)
	}
}
