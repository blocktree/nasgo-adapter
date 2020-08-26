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

package openwtester

import (
	"testing"

	"github.com/blocktree/openwallet/v2/openw"

	"github.com/blocktree/openwallet/v2/log"
	"github.com/blocktree/openwallet/v2/openwallet"
)

func testGetAssetsAccountBalance(tm *openw.WalletManager, walletID, accountID string) {
	balance, err := tm.GetAssetsAccountBalance(testApp, walletID, accountID)
	if err != nil {
		log.Error("GetAssetsAccountBalance failed, unexpected error:", err)
		return
	}
	log.Info("balance:", balance)
}

func testGetAssetsAccountTokenBalance(tm *openw.WalletManager, walletID, accountID string, contract openwallet.SmartContract) {
	balance, err := tm.GetAssetsAccountTokenBalance(testApp, walletID, accountID, contract)
	if err != nil {
		log.Error("GetAssetsAccountTokenBalance failed, unexpected error:", err)
		return
	}
	log.Info("token balance:", balance.Balance)
}

func testCreateTransactionStep(tm *openw.WalletManager, walletID, accountID, to, amount, feeRate string, contract *openwallet.SmartContract) (*openwallet.RawTransaction, error) {

	//err := tm.RefreshAssetsAccountBalance(testApp, accountID)
	//if err != nil {
	//	log.Error("RefreshAssetsAccountBalance failed, unexpected error:", err)
	//	return nil, err
	//}

	rawTx, err := tm.CreateTransaction(testApp, walletID, accountID, amount, to, feeRate, "", contract, nil)

	if err != nil {
		log.Error("CreateTransaction failed, unexpected error:", err)
		return nil, err
	}

	return rawTx, nil
}

func testCreateSummaryTransactionStep(
	tm *openw.WalletManager,
	walletID, accountID, summaryAddress, minTransfer, retainedBalance, feeRate string,
	start, limit int,
	contract *openwallet.SmartContract,
	feeSupportAccount *openwallet.FeesSupportAccount) ([]*openwallet.RawTransactionWithError, error) {

	rawTxArray, err := tm.CreateSummaryRawTransactionWithError(testApp, walletID, accountID, summaryAddress, minTransfer,
		retainedBalance, feeRate, start, limit, contract, feeSupportAccount)

	if err != nil {
		log.Error("CreateSummaryTransaction failed, unexpected error:", err)
		return nil, err
	}

	return rawTxArray, nil
}

func testSignTransactionStep(tm *openw.WalletManager, rawTx *openwallet.RawTransaction) (*openwallet.RawTransaction, error) {

	_, err := tm.SignTransaction(testApp, rawTx.Account.WalletID, rawTx.Account.AccountID, "12345678", rawTx)
	if err != nil {
		log.Error("SignTransaction failed, unexpected error:", err)
		return nil, err
	}

	log.Infof("rawTx: %+v", rawTx)
	return rawTx, nil
}

func testVerifyTransactionStep(tm *openw.WalletManager, rawTx *openwallet.RawTransaction) (*openwallet.RawTransaction, error) {

	//log.Info("rawTx.Signatures:", rawTx.Signatures)

	_, err := tm.VerifyTransaction(testApp, rawTx.Account.WalletID, rawTx.Account.AccountID, rawTx)
	if err != nil {
		log.Error("VerifyTransaction failed, unexpected error:", err)
		return nil, err
	}

	log.Infof("rawTx: %+v", rawTx)
	return rawTx, nil
}

func testSubmitTransactionStep(tm *openw.WalletManager, rawTx *openwallet.RawTransaction) (*openwallet.RawTransaction, error) {

	tx, err := tm.SubmitTransaction(testApp, rawTx.Account.WalletID, rawTx.Account.AccountID, rawTx)
	if err != nil {
		log.Error("SubmitTransaction failed, unexpected error:", err)
		return nil, err
	}

	log.Std.Info("tx: %+v", tx)
	log.Info("wxID:", tx.WxID)
	log.Info("txID:", rawTx.TxID)

	return rawTx, nil
}

func TestTransfer(t *testing.T) {

	tm := testInitWalletManager()
	walletID := "W2DyYXbPCpkXWS1tJPYcRxhioSNyqwSu8F"
	accountID := "9YBe43SkTyBneYNEnR7tHB3dh7VPB7toYkaZzU869C9y"
	// to := "NDt9qnAHnFAuP8T9GbzQ2o8UaacQscAcU2"
	//WMGcsvAwjjBj587oGE2GCZ3gu7F942hwGK
	//EhXYgY4wFN91VzkmJtyXPa1mPwEcp7o7PokQqaKcKGE4
	// to := "N2MzN3J9ZhHiWdmKxGSCxbwRHgWN7FzPC3"
	to := "NNd5jNQQ7E4p1s3QnUGgDTyZRaSb5asVkT"
	testGetAssetsAccountBalance(tm, walletID, accountID)

	// rawTx, err := testCreateTransactionStep(tm, walletID, accountID, to, "10", "", nil)
	rawTx, err := testCreateTransactionStep(tm, walletID, accountID, to, "0.1", "", nil)
	if err != nil {
		return
	}

	_, err = testSignTransactionStep(tm, rawTx)
	if err != nil {
		return
	}

	_, err = testVerifyTransactionStep(tm, rawTx)
	if err != nil {
		return
	}

	_, err = testSubmitTransactionStep(tm, rawTx)
	if err != nil {
		return
	}

}

func TestTransfer_Token(t *testing.T) {

	addrs := []string{
		"N2MzN3J9ZhHiWdmKxGSCxbwRHgWN7FzPC3",
		//"N4yAAjaTpW4Dc1c9Jd3k8tpJE2Q3ydQr4F",
		//"NC3VtNUNv2FR15yjrK5jhpbzLrY37LpaBg",
		//"NENz4Y4cMMnnAtsDmBneetMEupqZk4PGPt",
		//"NH9MCEEM1idr5JGVkPioEMXnQiCMwHHg1K",
		//"NJMrRD6BgfPTcJE8G7a5LKc4YMXEc5ARDa",

		//"NGJsqQJ9y8AMjGqXXZfLsS2NJdpMke3v6",
	}

	tm := testInitWalletManager()
	walletID := "W2DyYXbPCpkXWS1tJPYcRxhioSNyqwSu8F"
	accountID := "9YBe43SkTyBneYNEnR7tHB3dh7VPB7toYkaZzU869C9y"

	contract := openwallet.SmartContract{
		Address:  "IMM.IMM",
		Symbol:   "NSG",
		Name:     "IMMT",
		Token:    "IMMT",
		Decimals: 5,
	}

	testGetAssetsAccountBalance(tm, walletID, accountID)
	testGetAssetsAccountTokenBalance(tm, walletID, accountID, contract)

	for _, to := range addrs {
		// rawTx, err := testCreateTransactionStep(tm, walletID, accountID, to, "10", "", nil)
		rawTx, err := testCreateTransactionStep(tm, walletID, accountID, to, "0.1", "", &contract)
		if err != nil {
			return
		}

		_, err = testSignTransactionStep(tm, rawTx)
		if err != nil {
			return
		}

		_, err = testVerifyTransactionStep(tm, rawTx)
		if err != nil {
			return
		}

		_, err = testSubmitTransactionStep(tm, rawTx)
		if err != nil {
			return
		}
	}
}

func TestSummary(t *testing.T) {
	tm := testInitWalletManager()
	walletID := "WMGcsvAwjjBj587oGE2GCZ3gu7F942hwGK"
	accountID := "EhXYgY4wFN91VzkmJtyXPa1mPwEcp7o7PokQqaKcKGE4"
	summaryAddress := "N6E3HkfTUCpUA6F4RoDCEsNXzQ65HxJz3A"

	testGetAssetsAccountBalance(tm, walletID, accountID)

	feesSupport := &openwallet.FeesSupportAccount{
		AccountID:        "47VD3c4xUuvCu1cuaQffRMcgQdkkAtYovUwwiMNFpKNe", //NLYuCnWxigWcjJbmcwH6oKqH6zGGaoD9cc
		FeesSupportScale: "3",
	}

	rawTxArray, err := testCreateSummaryTransactionStep(tm, walletID, accountID,
		summaryAddress, "", "", "",
		0, 100, nil, feesSupport)
	// 0, 100, contract, nil)
	if err != nil {
		log.Errorf("CreateSummaryTransaction failed, unexpected error: %v", err)
		return
	}

	//执行汇总交易
	for _, rawTxWithErr := range rawTxArray {

		if rawTxWithErr.Error != nil {
			log.Error(rawTxWithErr.Error.Error())
			continue
		}

		_, err = testSignTransactionStep(tm, rawTxWithErr.RawTx)
		if err != nil {
			return
		}

		_, err = testVerifyTransactionStep(tm, rawTxWithErr.RawTx)
		if err != nil {
			return
		}

		_, err = testSubmitTransactionStep(tm, rawTxWithErr.RawTx)
		if err != nil {
			return
		}
	}

}

func TestSummary_Token(t *testing.T) {
	tm := testInitWalletManager()
	walletID := "WMGcsvAwjjBj587oGE2GCZ3gu7F942hwGK"
	accountID := "EhXYgY4wFN91VzkmJtyXPa1mPwEcp7o7PokQqaKcKGE4"
	summaryAddress := "N6E3HkfTUCpUA6F4RoDCEsNXzQ65HxJz3A"

	contract := openwallet.SmartContract{
		Address:  "IMM.IMM",
		Symbol:   "NSG",
		Name:     "IMMT",
		Token:    "IMMT",
		Decimals: 5,
	}

	testGetAssetsAccountBalance(tm, walletID, accountID)
	testGetAssetsAccountTokenBalance(tm, walletID, accountID, contract)

	feesSupport := &openwallet.FeesSupportAccount{
		AccountID:        "47VD3c4xUuvCu1cuaQffRMcgQdkkAtYovUwwiMNFpKNe", //NLYuCnWxigWcjJbmcwH6oKqH6zGGaoD9cc
		FeesSupportScale: "5",
	}

	rawTxArray, err := testCreateSummaryTransactionStep(tm, walletID, accountID,
		summaryAddress, "1", "0", "",
		0, 100, &contract, feesSupport)
	// 0, 100, contract, nil)
	if err != nil {
		log.Errorf("CreateSummaryTransaction failed, unexpected error: %v", err)
		return
	}

	//执行汇总交易
	for _, rawTxWithErr := range rawTxArray {

		if rawTxWithErr.Error != nil {
			log.Error(rawTxWithErr.Error.Error())
			continue
		}

		_, err = testSignTransactionStep(tm, rawTxWithErr.RawTx)
		if err != nil {
			return
		}

		_, err = testVerifyTransactionStep(tm, rawTxWithErr.RawTx)
		if err != nil {
			return
		}

		_, err = testSubmitTransactionStep(tm, rawTxWithErr.RawTx)
		if err != nil {
			return
		}
	}

}
