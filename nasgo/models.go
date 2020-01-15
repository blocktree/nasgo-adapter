/*
* Copyright 2020 The OpenWallet Authors
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
	"encoding/hex"
	"fmt"

	"github.com/assetsadapterstore/nasgo-adapter/rpc"
	"github.com/assetsadapterstore/nasgo-adapter/utils"
	"github.com/blocktree/go-owcrypt"
	"github.com/blocktree/openwallet/openwallet"
)

type Block struct {
	*rpc.Header
}

func (block *Block) BlockHeader() (header *openwallet.BlockHeader) {
	header = &openwallet.BlockHeader{}
	header.Hash = block.ID
	header.Version = uint64(block.Version)
	header.Time = uint64(block.Timestamp)
	header.Height = block.Height
	header.Previousblockhash = block.PrevBlock
	return
}

type Transaction struct {
	*rpc.Transaction
	SenderPublicKey []byte
}

func (tx *Transaction) GetID() string {
	hash := tx.GenerateHash()
	return hex.EncodeToString(hash)
}

func (tx *Transaction) GenerateHash() (hash []byte) {

	if tx == nil {
		fmt.Errorf("transaction is empty")
		return
	}
	if tx.Type != rpc.TxType_NSG || tx.Type != rpc.TxType_Asset {
		fmt.Errorf("transaction type is not allowed: %v", tx.Type)
		return
	}

	recipientId, _ := hex.DecodeString(tx.RecipientId)
	message, _ := hex.DecodeString(tx.Message)
	signature, _ := hex.DecodeString(tx.Signature)

	assetSlices := make([][]byte, 0)
	if tx.Type == rpc.TxType_Asset {
		if tx.Asset == nil {
			fmt.Errorf("transaction asset is empty")
			return
		}
		cur, _ := hex.DecodeString(tx.Asset.Currency)
		amt, _ := hex.DecodeString(tx.Asset.Amount)
		assetSlices = append(assetSlices, cur)
		assetSlices = append(assetSlices, amt)
	}
	assetSlice := utils.ConcatByteArray(assetSlices)

	txSlices := [][]byte{
		utils.UInt32ToBytes(tx.Type),
		utils.UInt64ToBytes(uint64(tx.Timestamp)),
		tx.SenderPublicKey,
		recipientId,
		utils.UInt64ToBytes(tx.Amount),
		message,
		assetSlice,
		signature,
	}

	msg := utils.ConcatByteArray(txSlices)
	hash = owcrypt.Hash(msg, 0, owcrypt.HASH_ALG_SHA256)
	return
}
