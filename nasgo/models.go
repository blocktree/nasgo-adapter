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
	"github.com/blocktree/nasgo-adapter/rpc"
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
