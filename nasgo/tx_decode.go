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
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/blocktree/nasgo-adapter/rpc"
	"github.com/blocktree/nasgo-adapter/txsigner"
	"github.com/blocktree/nasgo-adapter/utils"
	"github.com/blocktree/openwallet/openwallet"
	"github.com/shopspring/decimal"
)

type TransactionDecoder struct {
	openwallet.TransactionDecoderBase
	wm *WalletManager //钱包管理者
}

//NewTransactionDecoder 交易单解析器
func NewTransactionDecoder(wm *WalletManager) *TransactionDecoder {
	decoder := TransactionDecoder{}
	decoder.wm = wm
	return &decoder
}

//CreateRawTransaction 创建交易单
func (decoder *TransactionDecoder) CreateRawTransaction(wrapper openwallet.WalletDAI, rawTx *openwallet.RawTransaction) error {
	return decoder.CreateNSGRawTransaction(wrapper, rawTx)
}

//SignRawTransaction 签名交易单
func (decoder *TransactionDecoder) SignRawTransaction(wrapper openwallet.WalletDAI, rawTx *openwallet.RawTransaction) error {
	return decoder.SignNSGRawTransaction(wrapper, rawTx)
}

//VerifyRawTransaction 验证交易单，验证交易单并返回加入签名后的交易单
func (decoder *TransactionDecoder) VerifyRawTransaction(wrapper openwallet.WalletDAI, rawTx *openwallet.RawTransaction) error {
	return decoder.VerifyNSGRawTransaction(wrapper, rawTx)
}

// CreateSummaryRawTransactionWithError 创建汇总交易，返回能原始交易单数组（包含带错误的原始交易单）
func (decoder *TransactionDecoder) CreateSummaryRawTransactionWithError(wrapper openwallet.WalletDAI, sumRawTx *openwallet.SummaryRawTransaction) ([]*openwallet.RawTransactionWithError, error) {
	return nil, fmt.Errorf("do not support token transaction")
	// return decoder.CreateNSGSummaryRawTransaction(wrapper, sumRawTx)
}

//CreateSummaryRawTransaction 创建汇总交易，返回原始交易单数组
func (decoder *TransactionDecoder) CreateSummaryRawTransaction(wrapper openwallet.WalletDAI, sumRawTx *openwallet.SummaryRawTransaction) ([]*openwallet.RawTransaction, error) {
	var (
		rawTxWithErrArray []*openwallet.RawTransactionWithError
		rawTxArray        = make([]*openwallet.RawTransaction, 0)
		err               error
	)
	if sumRawTx.Coin.IsContract {
		return nil, fmt.Errorf("do not support token transaction")
	} else {
		// rawTxWithErrArray, err = decoder.CreateNSGSummaryRawTransaction(wrapper, sumRawTx)
	}
	if err != nil {
		return nil, err
	}
	for _, rawTxWithErr := range rawTxWithErrArray {
		if rawTxWithErr.Error != nil {
			continue
		}
		rawTxArray = append(rawTxArray, rawTxWithErr.RawTx)
	}
	return rawTxArray, nil
}

//SubmitRawTransaction 广播交易单
func (decoder *TransactionDecoder) SubmitRawTransaction(wrapper openwallet.WalletDAI, rawTx *openwallet.RawTransaction) (*openwallet.Transaction, error) {
	var trx txsigner.Transaction

	if len(rawTx.RawHex) == 0 {
		return nil, fmt.Errorf("transaction hex is empty")
	}

	if !rawTx.IsCompleted {
		return nil, fmt.Errorf("transaction is not completed validation")
	}

	rawHex, err := hex.DecodeString(rawTx.RawHex)
	if err != nil {
		return nil, openwallet.ConvertError(err)
	}

	err = json.Unmarshal(rawHex, &trx)
	if err != nil {
		return nil, openwallet.ConvertError(err)
	}

	param := map[string]interface{}{
		"transaction": trx,
	}

	err = decoder.wm.WalletClient.Tx.BroadcastTx(param)
	//err = decoder.wm.WalletClient.Tx.Broadcast(param)
	if err != nil {
		return nil, err
	}

	rawTx.TxID = trx.ID
	rawTx.IsSubmit = true

	decimals := int32(0)
	fees := rawTx.Fees
	if rawTx.Coin.IsContract {
		decimals = int32(rawTx.Coin.Contract.Decimals)
	} else {
		decimals = int32(decoder.wm.Decimal())
	}

	//记录一个交易单
	tx := &openwallet.Transaction{
		From:       rawTx.TxFrom,
		To:         rawTx.TxTo,
		Amount:     rawTx.TxAmount,
		Coin:       rawTx.Coin,
		TxID:       rawTx.TxID,
		Decimal:    decimals,
		AccountID:  rawTx.Account.AccountID,
		Fees:       fees,
		SubmitTime: time.Now().Unix(),
	}

	tx.WxID = openwallet.GenTransactionWxID(tx)

	return tx, nil
}

////////////////////////// NSG implement //////////////////////////

//CreateNSGRawTransaction 创建交易单
func (decoder *TransactionDecoder) CreateNSGRawTransaction(wrapper openwallet.WalletDAI, rawTx *openwallet.RawTransaction) error {

	var (
		balance   = decimal.New(0, 0)
		totalSend = decimal.New(0, 0)
		fixFees   = decimal.New(0, 0)
		from      = &openwallet.Address{}
		target    = ""
		accountID = rawTx.Account.AccountID
		limit     = 2000
		isToken   = rawTx.Coin.IsContract
		precision = int32(0)
	)

	if len(rawTx.To) == 0 {
		return errors.New("Receiver address is empty")
	}

	addresses, err := wrapper.GetAddressList(0, limit, "AccountID", rawTx.Account.AccountID)
	if err != nil {
		return err
	}

	if len(addresses) == 0 {
		return openwallet.Errorf(openwallet.ErrAccountNotAddress, "[%s] have not address", accountID)
	}

	//计算总发送金额
	for address, amount := range rawTx.To {
		amt, _ := decimal.NewFromString(amount)
		target = address
		totalSend = totalSend.Add(amt)
	}

	if len(rawTx.FeeRate) == 0 {
		fixFees, err = decimal.NewFromString(decoder.wm.Config.FixFees)
		if err != nil {
			return err
		}
	} else {
		fixFees, _ = decimal.NewFromString(rawTx.FeeRate)
	}

	decoder.wm.Log.Info("Calculating wallet unspent record to build transaction...")
	computeTotalSend := totalSend.Add(fixFees)

	//计算一个可用于支付的余额
	for _, addr := range addresses {
		if isToken {
			coin := rawTx.Coin.Contract.Address
			b, err := decoder.wm.WalletClient.Wallet.GetAssetsBalance(addr.Address, coin)
			if err != nil {
				return err
			}
			balance, _ = decimal.NewFromString(b.Balance)
			balance = balance.Shift(-int32(b.Precision))
			precision = int32(b.Precision)
		} else {
			b, err := decoder.wm.WalletClient.Wallet.GetBalance(addr.Address)
			if err != nil {
				return err
			}
			balance = decimal.New(int64(b), -decoder.wm.Decimal())
			precision = int32(decoder.wm.Decimal())
		}
		if balance.GreaterThanOrEqual(computeTotalSend) {
			from = addr
			break
		}
	}

	//判断余额是否足够支付发送数额+手续费
	if balance.LessThan(computeTotalSend) {
		return fmt.Errorf("The balance: %s is not enough! ", balance.StringFixed(decoder.wm.Decimal()))
	}

	rawTx.FeeRate = fixFees.StringFixed(decoder.wm.Decimal())
	rawTx.Fees = fixFees.StringFixed(decoder.wm.Decimal())

	decoder.wm.Log.Std.Notice("-----------------------------------------------")
	decoder.wm.Log.Std.Notice("From Account: %s", accountID)
	decoder.wm.Log.Std.Notice("To Address: %s", target)
	decoder.wm.Log.Std.Notice("Balance: %v", balance.String())
	decoder.wm.Log.Std.Notice("Fees: %v", fixFees.String())
	decoder.wm.Log.Std.Notice("Receive: %v", totalSend.String())
	decoder.wm.Log.Std.Notice("-----------------------------------------------")

	err = decoder.createNSGRawTransaction(wrapper, rawTx, from, target, totalSend.Shift(precision), fixFees)
	if err != nil {
		return err
	}

	return nil
}

//SignNSGRawTransaction 签名交易单
func (decoder *TransactionDecoder) SignNSGRawTransaction(wrapper openwallet.WalletDAI, rawTx *openwallet.RawTransaction) error {

	if rawTx.Signatures == nil || len(rawTx.Signatures) == 0 {
		return fmt.Errorf("transaction signature is empty")
	}

	key, err := wrapper.HDKey()
	if err != nil {
		return err
	}

	keySignatures := rawTx.Signatures[rawTx.Account.AccountID]
	if keySignatures != nil {
		for _, keySignature := range keySignatures {

			childKey, err := key.DerivedKeyWithPath(keySignature.Address.HDPath, keySignature.EccType)
			keyBytes, err := childKey.GetPrivateKeyBytes()
			if err != nil {
				return err
			}
			txHash := keySignature.Message
			decoder.wm.Log.Debug("hash:", txHash)

			data, err := hex.DecodeString(txHash)
			if err != nil {
				return fmt.Errorf("Invalid message to sign")
			}

			//签名交易
			/////////交易单哈希签名
			signature, err := txsigner.Default.SignTransactionHash(data, keyBytes, keySignature.EccType)
			if err != nil {
				return fmt.Errorf("transaction hash sign failed, unexpected error: %v", err)
			}

			keySignature.Signature = hex.EncodeToString(signature)
		}
	}

	decoder.wm.Log.Info("transaction hash sign success")

	rawTx.Signatures[rawTx.Account.AccountID] = keySignatures

	return nil
}

//VerifyNSGRawTransaction 验证交易单，验证交易单并返回加入签名后的交易单
func (decoder *TransactionDecoder) VerifyNSGRawTransaction(wrapper openwallet.WalletDAI, rawTx *openwallet.RawTransaction) error {

	rawHex, err := hex.DecodeString(rawTx.RawHex)
	if err != nil {
		return err
	}

	emptyTrans := string(rawHex)

	if rawTx.Signatures == nil || len(rawTx.Signatures) == 0 {
		return fmt.Errorf("transaction signature is empty")
	}

	for accountID, keySignatures := range rawTx.Signatures {
		decoder.wm.Log.Debug("accountID Signatures:", accountID)
		for _, keySignature := range keySignatures {

			signature, _ := hex.DecodeString(keySignature.Signature)
			pubkey, _ := hex.DecodeString(keySignature.Address.PublicKey)

			decoder.wm.Log.Debug("Message:", keySignature.Message)
			decoder.wm.Log.Debug("Signature:", keySignature.Signature)
			decoder.wm.Log.Debug("PublicKey:", keySignature.Address.PublicKey)

			msg, _ := hex.DecodeString(keySignature.Message)
			/////////验证交易单
			//验证时，对于公钥哈希地址，需要将对应的锁定脚本传入TxUnlock结构体
			pass, signedTrans, err := txsigner.Default.VerifyAndCombineTransaction(emptyTrans, msg, signature, pubkey)
			if pass {
				decoder.wm.Log.Debug("transaction verify passed")
				rawTx.IsCompleted = true
				rawTx.RawHex = signedTrans
			} else {
				decoder.wm.Log.Errorf("transaction verify failed, unexpected error: %v", err)
				rawTx.IsCompleted = false
			}
		}
	}

	return nil
}

//GetRawTransactionFeeRate 获取交易单的费率
func (decoder *TransactionDecoder) GetRawTransactionFeeRate() (feeRate string, unit string, err error) {
	return decoder.wm.Config.FixFees, "TX", nil
}

// //CreateNSGSummaryRawTransaction 创建汇总交易
// func (decoder *TransactionDecoder) CreateNSGSummaryRawTransaction(wrapper openwallet.WalletDAI, sumRawTx *openwallet.SummaryRawTransaction) ([]*openwallet.RawTransactionWithError, error) {

// 	var (
// 		accountID          = sumRawTx.Account.AccountID
// 		minTransfer, _     = decimal.NewFromString(sumRawTx.MinTransfer)
// 		retainedBalance, _ = decimal.NewFromString(sumRawTx.RetainedBalance)
// 		sumAddresses       = make([]string, 0)
// 		rawTxArray         = make([]*openwallet.RawTransactionWithError, 0)
// 		outputAddrs        map[string]decimal.Decimal
// 		totalInputAmount   decimal.Decimal
// 		sumUnspents        []*crypto.TransactionInputOutpoint
// 		fixFees            = decimal.New(0, 0)
// 	)

// 	if minTransfer.LessThan(retainedBalance) {
// 		return nil, fmt.Errorf("mini transfer amount must be greater than address retained balance")
// 	}

// 	address, err := wrapper.GetAddressList(sumRawTx.AddressStartIndex, sumRawTx.AddressLimit, "AccountID", sumRawTx.Account.AccountID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if len(address) == 0 {
// 		return nil, fmt.Errorf("[%s] have not addresses", accountID)
// 	}

// 	searchAddrs := make([]string, 0)
// 	for _, address := range address {
// 		searchAddrs = append(searchAddrs, address.Address)
// 	}

// 	addrBalanceArray, err := decoder.wm.Blockscanner.GetBalanceByAddress(searchAddrs...)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, addrBalance := range addrBalanceArray {
// 		decoder.wm.Log.Debugf("addrBalance: %+v", addrBalance)
// 		//检查余额是否超过最低转账
// 		addrBalanceDec, _ := decimal.NewFromString(addrBalance.Balance)
// 		if addrBalanceDec.GreaterThanOrEqual(minTransfer) {
// 			//添加到转账地址数组
// 			sumAddresses = append(sumAddresses, addrBalance.Address)
// 		}
// 	}

// 	if len(sumAddresses) == 0 {
// 		return nil, nil
// 	}

// 	//取得费率
// 	if len(sumRawTx.FeeRate) == 0 {
// 		fixFees, err = decimal.NewFromString(decoder.wm.Config.FixFees)
// 		if err != nil {
// 			return nil, err
// 		}
// 	} else {
// 		fixFees, _ = decimal.NewFromString(sumRawTx.FeeRate)
// 	}

// 	sumUnspents = make([]*crypto.TransactionInputOutpoint, 0)
// 	outputAddrs = make(map[string]decimal.Decimal, 0)
// 	totalInputAmount = decimal.Zero

// 	for i, addr := range sumAddresses {

// 		outputs, err := decoder.wm.WalletClient.Wallet.GetUnspent(addr)
// 		if err != nil {
// 			return nil, err
// 		}

// 		sumUnspents = append(sumUnspents, outputs...)
// 		if retainedBalance.GreaterThan(decimal.Zero) {
// 			outputAddrs = appendOutput(outputAddrs, addr, retainedBalance)
// 		}

// 		//如果遍历地址完结，就可以进行构建交易单
// 		if i == len(sumAddresses)-1 {

// 			//计算这笔交易单的汇总数量
// 			for _, u := range sumUnspents {
// 				ua := common.IntToDecimals(int64(u.Value), decoder.wm.Decimal())
// 				totalInputAmount = totalInputAmount.Add(ua)
// 			}

// 			/*

// 				汇总数量计算：

// 				1. 输入总数量 = 合计账户地址的所有utxo
// 				2. 账户地址输出总数量 = 账户地址保留余额 * 地址数
// 				3. 汇总数量 = 输入总数量 - 账户地址输出总数量 - 手续费
// 			*/
// 			retainedBalanceTotal := retainedBalance.Mul(decimal.New(int64(len(outputAddrs)), 0))
// 			sumAmount := totalInputAmount.Sub(retainedBalanceTotal).Sub(fixFees)

// 			decoder.wm.Log.Debugf("totalInputAmount: %v", totalInputAmount)
// 			decoder.wm.Log.Debugf("retainedBalanceTotal: %v", retainedBalanceTotal)
// 			decoder.wm.Log.Debugf("fees: %v", fixFees)
// 			decoder.wm.Log.Debugf("sumAmount: %v", sumAmount)

// 			//最后填充汇总地址及汇总数量
// 			outputAddrs = appendOutput(outputAddrs, sumRawTx.SummaryAddress, sumAmount)

// 			raxTxTo := make(map[string]string, 0)
// 			for a, m := range outputAddrs {
// 				raxTxTo[a] = m.StringFixed(decoder.wm.Decimal())
// 			}

// 			//创建一笔交易单
// 			rawTx := &openwallet.RawTransaction{
// 				Coin:     sumRawTx.Coin,
// 				Account:  sumRawTx.Account,
// 				FeeRate:  sumRawTx.FeeRate,
// 				To:       raxTxTo,
// 				Fees:     fixFees.StringFixed(decoder.wm.Decimal()),
// 				Required: 1,
// 			}

// 			createErr := decoder.createNSGRawTransaction(wrapper, rawTx, sumUnspents, outputAddrs, "", fixFees)
// 			rawTxWithErr := &openwallet.RawTransactionWithError{
// 				RawTx: rawTx,
// 				Error: openwallet.ConvertError(createErr),
// 			}

// 			//创建成功，添加到队列
// 			rawTxArray = append(rawTxArray, rawTxWithErr)

// 			//清空临时变量
// 			sumUnspents = make([]*crypto.TransactionInputOutpoint, 0)
// 			outputAddrs = make(map[string]decimal.Decimal, 0)
// 			totalInputAmount = decimal.Zero

// 		}
// 	}

// 	return rawTxArray, nil
// }

//createNSGRawTransaction 创建NSG原始交易单
func (decoder *TransactionDecoder) createNSGRawTransaction(
	wrapper openwallet.WalletDAI,
	rawTx *openwallet.RawTransaction,
	from *openwallet.Address,
	to string,
	amount decimal.Decimal,
	fees decimal.Decimal,
) error {

	var (
		accountTotalSent = decimal.Zero
	)

	if len(to) == 0 {
		return fmt.Errorf("Receiver addresses is empty! ")
	}
	accountTotalSent = accountTotalSent.Add(amount)

	trx := &txsigner.Transaction{}
	trx.Transaction = &rpc.Transaction{}
	if rawTx.Coin.IsContract {
		trx.Asset.Currency = rawTx.Coin.Contract.Address
		trx.Asset.Amount = amount.String()
		trx.Type = rpc.TxType_Asset
	} else {
		trx.Amount = uint64(amount.IntPart())
		trx.Type = rpc.TxType_NSG
	}
	trx.Timestamp = utils.GetEpochTime()
	trx.SenderPublicKey = from.PublicKey
	trx.RecipientId = to
	trx.Message = rawTx.GetExtParam().Get("memo").String()

	//trx.ID = trx.GetID()
	txBytes, err := json.Marshal(trx)
	if err != nil {
		return openwallet.Errorf(openwallet.ErrCreateRawTransactionFailed, "Failed to marshal transaction: %s", err)
	}
	rawTx.RawHex = hex.EncodeToString(txBytes)

	if rawTx.Signatures == nil {
		rawTx.Signatures = make(map[string][]*openwallet.KeySignature)
	}
	//装配签名
	keySigs := make([]*openwallet.KeySignature, 0)
	trxHash := trx.GenerateHash(true)
	beSignHex := hex.EncodeToString(trxHash)

	decoder.wm.Log.Std.Debug("txHash: %s", beSignHex)

	signature := openwallet.KeySignature{
		EccType: decoder.wm.Config.CurveType,
		Nonce:   "",
		Address: from,
		Message: beSignHex,
	}

	keySigs = append(keySigs, &signature)

	feesDec, _ := decimal.NewFromString(rawTx.Fees)
	accountTotalSent = accountTotalSent.Add(feesDec)
	accountTotalSent = decimal.Zero.Sub(accountTotalSent)

	rawTx.Signatures[rawTx.Account.AccountID] = keySigs
	rawTx.IsBuilt = true
	rawTx.TxAmount = accountTotalSent.StringFixed(decoder.wm.Decimal())
	rawTx.TxFrom = []string{from.Address}
	rawTx.TxTo = []string{to}

	return nil
}
