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
	"fmt"
	"strconv"
	"time"

	"github.com/blocktree/nasgo-adapter/rpc"
	"github.com/blocktree/openwallet/log"
	"github.com/blocktree/openwallet/openwallet"
	"github.com/shopspring/decimal"
)

const (
	blockchainBucket = "blockchain" // blockchain dataset
	//periodOfTask      = 5 * time.Second // task interval
	maxExtractingSize = 10 // thread count
)

//BlockScanner block scanner
type BlockScanner struct {
	*openwallet.BlockScannerBase

	CurrentBlockHeight   uint64         //当前区块高度
	extractingCH         chan struct{}  //扫描工作令牌
	wm                   *WalletManager //钱包管理者
	IsScanMemPool        bool           //是否扫描交易池
	RescanLastBlockCount uint64         //重扫上N个区块数量
}

//ExtractResult extract result
type ExtractResult struct {
	extractData map[string][]*openwallet.TxExtractData
	TxID        string
	BlockHash   string
	BlockHeight uint64
	BlockTime   int64
	Success     bool
}

//SaveResult result
type SaveResult struct {
	TxID        string
	BlockHeight uint64
	Success     bool
}

// NewBlockScanner create a block scanner
func NewBlockScanner(wm *WalletManager) *BlockScanner {
	bs := BlockScanner{
		BlockScannerBase: openwallet.NewBlockScannerBase(),
	}

	bs.extractingCH = make(chan struct{}, maxExtractingSize)
	bs.wm = wm
	bs.IsScanMemPool = false
	bs.RescanLastBlockCount = 0

	// set task
	bs.SetTask(bs.ScanBlockTask)

	return &bs
}

// ScanBlockTask scan block task
func (bs *BlockScanner) ScanBlockTask() {

	var (
		currentHeight uint32
		currentHash   string
	)

	// get local block header
	currentHeight, currentHash, err := bs.GetLocalBlockHead()

	if err != nil {
		bs.wm.Log.Std.Error("", err)
	}

	if currentHeight == 0 {
		bs.wm.Log.Std.Info("No records found in local, get current block as the local!")

		headBlock, err := bs.GetGlobalHeadBlock()
		if err != nil {
			bs.wm.Log.Std.Info("get head block error, err=%v", err)
		}

		currentHash = headBlock.Header.PrevBlock
		currentHeight = uint32(headBlock.Height - 1)
	}

	for {
		if !bs.Scanning {
			// stop scan
			return
		}

		maxBlockHeight := bs.GetGlobalMaxBlockHeight()

		bs.wm.Log.Info("current block height:", currentHeight, " maxBlockHeight:", maxBlockHeight)
		if uint64(currentHeight) == maxBlockHeight-1 {
			bs.wm.Log.Std.Info("block scanner has scanned full chain data. Current height %d", maxBlockHeight)
			break
		}

		// next block
		currentHeight = currentHeight + 1

		bs.wm.Log.Std.Info("block scanner scanning height: %d ...", currentHeight)
		block, err := bs.GetByHeight(currentHeight)

		if err != nil {
			bs.wm.Log.Std.Info("block scanner can not get new block data by rpc; unexpected error: %v", err)
			break
		}

		if currentHash != block.PrevBlock {
			bs.wm.Log.Std.Info("block has been fork on height: %d.", currentHeight)
			bs.wm.Log.Std.Info("block height: %d local hash = %s ", currentHeight-1, currentHash)
			bs.wm.Log.Std.Info("block height: %d mainnet hash = %s ", currentHeight-1, block.PrevBlock)
			bs.wm.Log.Std.Info("delete recharge records on block height: %d.", currentHeight-1)

			// get local fork bolck
			forkBlock, _ := bs.GetLocalBlock(currentHeight - 1)
			// delete last unscan block
			bs.DeleteUnscanRecord(currentHeight - 1)
			currentHeight = currentHeight - 2 // scan back to last 2 block
			if currentHeight <= 0 {
				currentHeight = 1
			}
			localBlock, err := bs.GetLocalBlock(currentHeight)
			if err != nil {
				bs.wm.Log.Std.Error("block scanner can not get local block; unexpected error: %v", err)
				//get block from rpc
				bs.wm.Log.Info("block scanner prev block height:", currentHeight)
				curBlock, err := bs.GetByHeight(currentHeight)
				if err != nil {
					bs.wm.Log.Std.Error("block scanner can not get prev block by rpc; unexpected error: %v", err)
					break
				}
				currentHash = curBlock.ID
			} else {
				//重置当前区块的hash
				currentHash = localBlock.ID
			}
			bs.wm.Log.Std.Info("rescan block on height: %d, hash: %s .", currentHeight, currentHash)

			//重新记录一个新扫描起点
			bs.SaveLocalBlockHead(currentHeight, currentHash)

			if forkBlock != nil {
				//通知分叉区块给观测者，异步处理
				bs.forkBlockNotify(forkBlock)
			}

		} else {
			currentHash = block.ID
			err := bs.BatchExtractTransactions(uint64(currentHeight), currentHash, block.Timestamp)
			if err != nil {
				bs.wm.Log.Std.Error("block scanner ran BatchExtractTransactions occured unexpected error: %v", err)
			}

			//保存本地新高度
			bs.SaveLocalBlockHead(currentHeight, currentHash)
			bs.SaveLocalBlock(block)
			//通知新区块给观测者，异步处理
			bs.newBlockNotify(block)
		}
	}

	//重扫失败区块
	bs.RescanFailedRecord()

}

//newBlockNotify 获得新区块后，通知给观测者
func (bs *BlockScanner) forkBlockNotify(block *Block) {
	header := block.BlockHeader(bs.wm.Symbol())
	header.Fork = true
	bs.NewBlockNotify(header)
}

//newBlockNotify 获得新区块后，通知给观测者
func (bs *BlockScanner) newBlockNotify(block *Block) {
	header := block.BlockHeader(bs.wm.Symbol())
	bs.NewBlockNotify(header)
}

// BatchExtractTransactions 批量提取交易单
func (bs *BlockScanner) BatchExtractTransactions(blockHeight uint64, blockHash string, blockTime int64) error {

	var (
		quit       = make(chan struct{})
		done       = 0 //完成标记
		failed     = 0
		shouldDone = 0 //需要完成的总数
	)

	transactions, err := bs.wm.WalletClient.Tx.GetTransactionsByBlock(blockHash)
	if err != nil {
		bs.wm.Log.Std.Info("block scanner can not get transactions; unexpected error: %v", err)
	}

	if len(transactions) == 0 {
		return nil
	}

	shouldDone = len(transactions)
	bs.wm.Log.Std.Info("block scanner ready extract transactions total: %d ", len(transactions))

	//生产通道
	producer := make(chan ExtractResult)
	defer close(producer)

	//消费通道
	worker := make(chan ExtractResult)
	defer close(worker)

	//保存工作
	saveWork := func(height uint64, result chan ExtractResult) {
		//回收创建的地址
		for gets := range result {

			if gets.Success {
				notifyErr := bs.newExtractDataNotify(height, gets.extractData)
				if notifyErr != nil {
					failed++ //标记保存失败数
					bs.wm.Log.Std.Info("newExtractDataNotify unexpected error: %v", notifyErr)
				}
			} else {
				//记录未扫区块
				unscanRecord := openwallet.NewUnscanRecord(height, "", "", bs.wm.Symbol())
				bs.SaveUnscanRecord(unscanRecord)
				failed++ //标记保存失败数
			}
			//累计完成的线程数
			done++
			if done == shouldDone {
				close(quit) //关闭通道，等于给通道传入nil
			}
		}
	}

	//提取工作
	extractWork := func(eblockHeight uint64, eBlockHash string, eBlockTime int64, mTransactions []*rpc.Transaction, eProducer chan ExtractResult) {
		for _, tx := range mTransactions {
			bs.extractingCH <- struct{}{}

			go func(mBlockHeight uint64, mTx *rpc.Transaction, end chan struct{}, mProducer chan<- ExtractResult) {
				//导出提出的交易
				mProducer <- bs.ExtractTransaction(mBlockHeight, eBlockHash, eBlockTime, mTx, bs.ScanTargetFunc)
				//释放
				<-end

			}(eblockHeight, tx, bs.extractingCH, eProducer)
		}
	}
	/*	开启导出的线程	*/

	//独立线程运行消费
	go saveWork(blockHeight, worker)

	//独立线程运行生产
	go extractWork(blockHeight, blockHash, blockTime, transactions, producer)

	//以下使用生产消费模式
	bs.extractRuntime(producer, worker, quit)

	if failed > 0 {
		return fmt.Errorf("block scanner saveWork failed")
	}

	return nil
}

//extractRuntime 提取运行时
func (bs *BlockScanner) extractRuntime(producer chan ExtractResult, worker chan ExtractResult, quit chan struct{}) {

	var (
		values = make([]ExtractResult, 0)
	)

	for {
		var activeWorker chan<- ExtractResult
		var activeValue ExtractResult
		//当数据队列有数据时，释放顶部，传输给消费者
		if len(values) > 0 {
			activeWorker = worker
			activeValue = values[0]
		}
		select {
		//生成者不断生成数据，插入到数据队列尾部
		case pa := <-producer:
			values = append(values, pa)
		case <-quit:
			//退出
			return
		case activeWorker <- activeValue:
			values = values[1:]
		}
	}
	//return
}

// ExtractTransaction 提取交易单
func (bs *BlockScanner) ExtractTransaction(blockHeight uint64, blockHash string, blockTime int64, trx *rpc.Transaction, scanTargetFunc openwallet.BlockScanTargetFunc) ExtractResult {
	var (
		success = true
		result  = ExtractResult{
			BlockHash:   blockHash,
			BlockHeight: blockHeight,
			TxID:        trx.ID,
			extractData: make(map[string][]*openwallet.TxExtractData),
			BlockTime:   blockTime,
		}
		err error
	)

	if trx.Type == rpc.TxType_Asset {
		txid := trx.ID
		trx, err = bs.wm.WalletClient.Tx.GetTransaction(txid)
		if err != nil {
			bs.wm.Log.Std.Debug("get asset transaction fail: [%v] ", txid)
			return ExtractResult{Success: true}
		}
	} else if trx.Type != 0 {
		bs.wm.Log.Std.Debug("does not support transaction type: [%v] ", trx.Type)
		return ExtractResult{Success: true}
	}

	if scanTargetFunc == nil {
		bs.wm.Log.Std.Error("scanTargetFunc is not configurated")
		return ExtractResult{Success: false}
	}

	from := trx.SenderID
	to := trx.RecipientId

	//订阅地址为交易单中的发送者
	accountID1, ok1 := scanTargetFunc(openwallet.ScanTarget{Address: from, Symbol: bs.wm.Symbol(), BalanceModelType: openwallet.BalanceModelTypeAddress})
	//订阅地址为交易单中的接收者
	accountID2, ok2 := scanTargetFunc(openwallet.ScanTarget{Address: to, Symbol: bs.wm.Symbol(), BalanceModelType: openwallet.BalanceModelTypeAddress})
	if accountID1 == accountID2 && len(accountID1) > 0 && len(accountID2) > 0 {
		bs.InitExtractResult(accountID1, trx, &result, 0)
	} else {
		if ok1 {
			bs.InitExtractResult(accountID1, trx, &result, 1)
		}

		if ok2 {
			bs.InitExtractResult(accountID2, trx, &result, 2)
		}
	}

	result.Success = success
	return result

}

//InitExtractResult optType = 0: 输入输出提取，1: 输入提取，2：输出提取
func (bs *BlockScanner) InitExtractResult(sourceKey string, trx *rpc.Transaction, result *ExtractResult, optType int64) {

	txExtractDataArray := result.extractData[sourceKey]
	if txExtractDataArray == nil {
		txExtractDataArray = make([]*openwallet.TxExtractData, 0)
	}

	txExtractData := &openwallet.TxExtractData{}

	status := "1"
	reason := ""
	amount := decimal.New(int64(trx.Amount), -bs.wm.Decimal()).String()
	from := trx.SenderID
	to := trx.RecipientId
	coin := openwallet.Coin{
		Symbol:     bs.wm.Symbol(),
		IsContract: trx.Type == rpc.TxType_Asset,
	}

	if trx.Type == rpc.TxType_Asset {
		if trx.Asset == nil || trx.Asset.UiaTransfer == nil {
			bs.wm.Log.Std.Debug("transaction asset info missing: [%v] ", trx.ID)
			return
		}
		token := trx.Asset.UiaTransfer.Currency
		contractID := openwallet.GenContractID(bs.wm.Symbol(), token)
		coin.Contract = openwallet.SmartContract{
			Symbol:     bs.wm.Symbol(),
			ContractID: contractID,
			Address:    token,
			Decimals:   uint64(trx.Asset.UiaTransfer.Precision),
		}
		amount = trx.Asset.UiaTransfer.Amount

		if optType == 0 || optType == 1 {
			fees := decimal.New(int64(trx.Fee), -bs.wm.Decimal()).String()
			feeExtractData := &openwallet.TxExtractData{}
			feeTransx := &openwallet.Transaction{
				Fees:        fees,
				BlockHash:   result.BlockHash,
				BlockHeight: result.BlockHeight,
				TxID:        result.TxID,
				Amount:      "0",
				ConfirmTime: result.BlockTime,
				From:        []string{from + ":" + amount},
				To:          []string{"'' :" + amount},
				IsMemo:      false,
				Status:      status,
				Reason:      reason,
				TxType:      1,
			}

			wxID := openwallet.GenTransactionWxID(feeTransx)
			feeTransx.WxID = wxID
			feeExtractData.Transaction = feeTransx

			feeCharge := &openwallet.TxInput{}
			feeCharge.Amount = fees
			feeCharge.TxType = feeTransx.TxType
			feeExtractData.TxInputs = append(feeExtractData.TxInputs, feeCharge)

			txExtractDataArray = append(txExtractDataArray, feeExtractData)
		}
	}

	transx := &openwallet.Transaction{
		Coin:        coin,
		BlockHash:   result.BlockHash,
		BlockHeight: result.BlockHeight,
		TxID:        result.TxID,
		Fees:        "0",
		Amount:      amount,
		ConfirmTime: result.BlockTime,
		From:        []string{from + ":" + amount},
		To:          []string{to + ":" + amount},
		IsMemo:      true,
		Status:      status,
		Reason:      reason,
		TxType:      0,
	}
	if trx.Type == rpc.TxType_NSG {
		transx.Fees = decimal.New(int64(trx.Fee), -bs.wm.Decimal()).String()
	}

	transx.SetExtParam("memo", trx.Message)

	wxID := openwallet.GenTransactionWxID(transx)
	transx.WxID = wxID

	txExtractData.Transaction = transx
	if optType == 0 {
		bs.extractTxInput(trx, txExtractData)
		bs.extractTxOutput(trx, txExtractData)
	} else if optType == 1 {
		bs.extractTxInput(trx, txExtractData)
	} else if optType == 2 {
		bs.extractTxOutput(trx, txExtractData)
	}

	txExtractDataArray = append(txExtractDataArray, txExtractData)

	result.extractData[sourceKey] = txExtractDataArray
}

//extractTxInput 提取交易单输入部分,无需手续费，所以只包含1个TxInput
func (bs *BlockScanner) extractTxInput(trx *rpc.Transaction, txExtractData *openwallet.TxExtractData) {

	tx := txExtractData.Transaction
	coin := openwallet.Coin(tx.Coin)

	from := trx.SenderID

	//主网from交易转账信息，第一个TxInput
	txInput := &openwallet.TxInput{}
	txInput.Recharge.Sid = openwallet.GenTxInputSID(tx.TxID, bs.wm.Symbol(), coin.ContractID, uint64(0))
	txInput.Recharge.TxID = tx.TxID
	txInput.Recharge.Address = from
	txInput.Recharge.Coin = coin
	txInput.Recharge.Amount = tx.Amount
	txInput.Recharge.Symbol = coin.Symbol
	txInput.Recharge.BlockHash = tx.BlockHash
	txInput.Recharge.BlockHeight = tx.BlockHeight
	txInput.Recharge.Index = 0
	txInput.Recharge.CreateAt = time.Now().Unix()
	txInput.Recharge.TxType = tx.TxType
	txExtractData.TxInputs = append(txExtractData.TxInputs, txInput)

	if trx.Type == rpc.TxType_NSG && trx.Fee > 0 {
		//手续费也作为一个输出s
		fees := decimal.New(int64(trx.Fee), -bs.wm.Decimal()).String()
		tmp := *txInput
		feeCharge := &tmp
		feeCharge.Amount = fees
		feeCharge.TxType = tx.TxType
		txExtractData.TxInputs = append(txExtractData.TxInputs, feeCharge)
	}

}

//extractTxOutput 提取交易单输入部分,只有一个TxOutPut
func (bs *BlockScanner) extractTxOutput(trx *rpc.Transaction, txExtractData *openwallet.TxExtractData) {

	tx := txExtractData.Transaction
	coin := openwallet.Coin(tx.Coin)
	to := trx.RecipientId

	//主网to交易转账信息,只有一个TxOutPut
	txOutput := &openwallet.TxOutPut{}
	txOutput.Recharge.Sid = openwallet.GenTxOutPutSID(tx.TxID, bs.wm.Symbol(), coin.ContractID, uint64(0))
	txOutput.Recharge.TxID = tx.TxID
	txOutput.Recharge.Address = to
	txOutput.Recharge.Coin = coin
	txOutput.Recharge.Amount = tx.Amount
	txOutput.Recharge.Symbol = coin.Symbol
	txOutput.Recharge.BlockHash = tx.BlockHash
	txOutput.Recharge.BlockHeight = tx.BlockHeight
	txOutput.Recharge.Index = 0
	txOutput.Recharge.CreateAt = time.Now().Unix()
	txExtractData.TxOutputs = append(txExtractData.TxOutputs, txOutput)
}

//newExtractDataNotify 发送通知
func (bs *BlockScanner) newExtractDataNotify(height uint64, extractData map[string][]*openwallet.TxExtractData) error {
	for o := range bs.Observers {
		for key, array := range extractData {
			for _, item := range array {
				err := o.BlockExtractDataNotify(key, item)
				if err != nil {
					log.Error("BlockExtractDataNotify unexpected error:", err)
					//记录未扫区块
					unscanRecord := openwallet.NewUnscanRecord(height, "", "ExtractData Notify failed.", bs.wm.Symbol())
					err = bs.SaveUnscanRecord(unscanRecord)
					if err != nil {
						log.Std.Error("block height: %d, save unscan record failed. unexpected error: %v", height, err.Error())
					}

				}
			}

		}
	}

	return nil
}

//ScanBlock 扫描指定高度区块
func (bs *BlockScanner) ScanBlock(height uint64) error {

	block, err := bs.scanBlock(height)
	if err != nil {
		return err
	}

	//通知新区块给观测者，异步处理
	bs.newBlockNotify(block)

	return nil
}

func (bs *BlockScanner) scanBlock(height uint64) (*Block, error) {

	block, err := bs.wm.WalletClient.Block.GetByHeight(height)
	if err != nil {
		bs.wm.Log.Std.Info("block scanner can not get new block data; unexpected error: %v", err)

		//记录未扫区块
		unscanRecord := openwallet.NewUnscanRecord(height, "", err.Error(), bs.wm.Symbol())
		bs.SaveUnscanRecord(unscanRecord)
		bs.wm.Log.Std.Info("block height: %d extract failed.", height)
		return nil, err
	}

	bs.wm.Log.Std.Info("block scanner scanning height: %d ...", block.ID)

	err = bs.BatchExtractTransactions(block.Height, block.ID, block.Timestamp)
	if err != nil {
		bs.wm.Log.Std.Info("block scanner can not extractRechargeRecords; unexpected error: %v", err)
	}

	return &Block{block}, nil
}

//SetRescanBlockHeight 重置区块链扫描高度
func (bs *BlockScanner) SetRescanBlockHeight(height uint64) error {
	if height <= 0 {
		return fmt.Errorf("block height to rescan must greater than 0. ")
	}

	block, err := bs.wm.WalletClient.Block.GetByHeight(height - 1)
	if err != nil {
		return err
	}

	bs.SaveLocalBlockHead(uint32(height-1), block.ID)

	return nil
}

// GetGlobalMaxBlockHeight GetGlobalMaxBlockHeight
func (bs *BlockScanner) GetGlobalMaxBlockHeight() uint64 {
	height, err := bs.wm.WalletClient.Block.GetBlockHeight()
	if err != nil {
		bs.wm.Log.Std.Info("block scanner can not get height; unexpected error:%v", err)
		return 0
	}
	return height
}

//GetGlobalHeadBlock GetGlobalHeadBlock
func (bs *BlockScanner) GetGlobalHeadBlock() (block *Block, err error) {

	height := bs.GetGlobalMaxBlockHeight()
	if height == 0 {
		bs.wm.Log.Std.Info("block scanner can not get height; unexpected error:%v", err)
		return
	}

	header, err := bs.wm.WalletClient.Block.GetByHeight(height - 1)
	if err != nil {
		bs.wm.Log.Std.Info("block scanner can not get block by height; unexpected error:%v", err)
		return
	}

	block = &Block{header}

	return
}

// GetByHeight GetByHeight
func (bs *BlockScanner) GetByHeight(height uint32) (block *Block, err error) {
	header, err := bs.wm.WalletClient.Block.GetByHeight(uint64(height))

	if err != nil {
		return nil, fmt.Errorf("block scanner can not get new block data by rpc; unexpected error: %v", err)
	}
	block = &Block{header}
	return
}

//GetScannedBlockHeight 获取已扫区块高度
func (bs *BlockScanner) GetScannedBlockHeight() uint64 {
	height, _, _ := bs.GetLocalBlockHead()
	return uint64(height)
}

//GetBalanceByAddress 查询地址余额
func (bs *BlockScanner) GetBalanceByAddress(address ...string) ([]*openwallet.Balance, error) {

	addrBalanceArr := make([]*openwallet.Balance, 0)

	for _, addr := range address {

		balance, err := bs.wm.WalletClient.Wallet.GetBalance(addr)
		if err != nil {
			bs.wm.Log.Errorf("get account[%v] token balance failed, err: %v", addr, err)
		}

		value := decimal.New(int64(balance), -bs.wm.Decimal())

		tokenBalance := &openwallet.Balance{
			Address:          addr,
			Symbol:           bs.wm.Symbol(),
			Balance:          value.String(),
			ConfirmBalance:   value.String(),
			UnconfirmBalance: "0",
		}

		addrBalanceArr = append(addrBalanceArr, tokenBalance)
	}

	return addrBalanceArr, nil
}

func (bs *BlockScanner) GetCurrentBlockHeader() (*openwallet.BlockHeader, error) {
	block, err := bs.GetGlobalHeadBlock()
	if err != nil {
		bs.wm.Log.Std.Info("get chain info error;unexpected error:%v", err)
		return nil, err
	}
	return block.BlockHeader(bs.wm.Symbol()), nil
}

//rescanFailedRecord 重扫失败记录
func (bs *BlockScanner) RescanFailedRecord() {

	var (
		blockMap = make(map[uint64][]string)
	)

	list, err := bs.BlockchainDAI.GetUnscanRecords(bs.wm.Symbol())
	if err != nil {
		bs.wm.Log.Std.Info("block scanner can not get rescan data; unexpected error: %v", err)
	}

	//组合成批处理
	for _, r := range list {

		if _, exist := blockMap[r.BlockHeight]; !exist {
			blockMap[r.BlockHeight] = make([]string, 0)
		}

		if len(r.TxID) > 0 {
			arr := blockMap[r.BlockHeight]
			arr = append(arr, r.TxID)

			blockMap[r.BlockHeight] = arr
		}
	}

	for height, _ := range blockMap {

		if height == 0 {
			continue
		}

		bs.wm.Log.Std.Info("block scanner rescanning height: %d ...", height)

		block, err := bs.wm.WalletClient.Block.GetByHeight(height)
		if err != nil {
			bs.wm.Log.Std.Info("block scanner can not get new block data; unexpected error: %v", err)
			continue
		}

		err = bs.BatchExtractTransactions(height, block.ID, block.Timestamp)
		if err != nil {
			bs.wm.Log.Std.Info("block scanner can not extractRechargeRecords; unexpected error: %v", err)
			continue
		}

		//删除未扫记录
		bs.DeleteUnscanRecord(uint32(height))
	}
}

//ExtractTransactionData 扫描一笔交易
func (bs *BlockScanner) ExtractTransactionData(txid string, scanTargetFunc openwallet.BlockScanTargetFunc) (map[string][]*openwallet.TxExtractData, error) {

	tx, err := bs.wm.WalletClient.Tx.GetTransaction(txid)
	if err != nil {
		return nil, err
	}

	height, err := strconv.ParseUint(tx.Height, 10, 64)
	block, err := bs.wm.WalletClient.Block.GetByHeight(height)

	result := bs.ExtractTransaction(block.Height, block.ID, block.Timestamp, tx, scanTargetFunc)
	return result.extractData, nil
}

//SupportBlockchainDAI 支持外部设置区块链数据访问接口
//@optional
func (bs *BlockScanner) SupportBlockchainDAI() bool {
	return true
}
