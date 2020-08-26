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
	"github.com/blocktree/nasgo-adapter/rpc"
	"github.com/blocktree/openwallet/v2/log"
	"github.com/blocktree/openwallet/v2/openwallet"
)

const (
	maxAddresNum = 10000
)

type WalletManager struct {
	openwallet.AssetsAdapterBase

	WalletClient    *rpc.Client                     // 节点客户端
	Config          *WalletConfig                   //钱包管理配置
	Blockscanner    *BlockScanner                   //区块扫描器
	Decoder         *AddressDecoder                 //地址编码器
	TxDecoder       openwallet.TransactionDecoder   //交易单编码器
	ContractDecoder openwallet.SmartContractDecoder //智能合约解析器
	Log             *log.OWLogger                   //日志工具
}

func NewWalletManager() *WalletManager {
	wm := WalletManager{}
	wm.Config = NewConfig(Symbol)
	wm.WalletClient = rpc.NewClient(wm.Config.ServerAPI)
	//区块扫描器
	wm.Blockscanner = NewBlockScanner(&wm)
	wm.Decoder = NewAddressDecoder(&wm)
	wm.TxDecoder = NewTransactionDecoder(&wm)
	wm.Log = log.NewOWLogger(wm.Symbol())
	wm.ContractDecoder = NewContractDecoder(&wm)
	return &wm
}
