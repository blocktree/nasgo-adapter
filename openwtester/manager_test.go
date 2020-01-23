package openwtester

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/blocktree/openwallet/log"
	"github.com/blocktree/openwallet/openw"
	"github.com/blocktree/openwallet/openwallet"
)

var (
	testApp        = "nasgo-adapter"
	configFilePath = filepath.Join("conf")
	dbFilePath     = filepath.Join("data", "db")
	dbFileName     = "blockchain-NSG.db"
)

func testInitWalletManager() *openw.WalletManager {
	log.SetLogFuncCall(true)
	tc := openw.NewConfig()

	tc.ConfigDir = configFilePath
	tc.EnableBlockScan = false
	tc.SupportAssets = []string{
		"NSG",
	}
	return openw.NewWalletManager(tc)
	//tm.Init()
}

func TestWalletManager_CreateWallet(t *testing.T) {
	tm := testInitWalletManager()
	w := &openwallet.Wallet{Alias: "HELLO NSG!!", IsTrust: true, Password: "12345678"}
	nw, key, err := tm.CreateWallet(testApp, w)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("wallet:", nw)
	log.Info("key:", key)

}

func TestWalletManager_GetWalletInfo(t *testing.T) {

	tm := testInitWalletManager()

	wallet, err := tm.GetWalletInfo(testApp, "W2DyYXbPCpkXWS1tJPYcRxhioSNyqwSu8F")
	if err != nil {
		log.Error("unexpected error:", err)
		return
	}
	log.Info("wallet:", wallet)
}

func TestWalletManager_GetWalletList(t *testing.T) {

	tm := testInitWalletManager()

	list, err := tm.GetWalletList(testApp, 0, 10000000)
	if err != nil {
		log.Error("unexpected error:", err)
		return
	}
	for i, w := range list {
		log.Info("wallet[", i, "] :", w)
	}
	log.Info("wallet count:", len(list))

	tm.CloseDB(testApp)
}

func TestWalletManager_CreateAssetsAccount(t *testing.T) {

	tm := testInitWalletManager()

	walletID := "W2DyYXbPCpkXWS1tJPYcRxhioSNyqwSu8F"
	account := &openwallet.AssetsAccount{Alias: "fee support NSG", WalletID: walletID, Required: 1, Symbol: "NSG", IsTrust: true}
	account, address, err := tm.CreateAssetsAccount(testApp, walletID, "12345678", account, nil)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("account:", account)
	log.Info("address:", address)

	tm.CloseDB(testApp)
}

func TestWalletManager_GetAssetsAccountList(t *testing.T) {

	tm := testInitWalletManager()

	walletID := "W2DyYXbPCpkXWS1tJPYcRxhioSNyqwSu8F"
	list, err := tm.GetAssetsAccountList(testApp, walletID, 0, 10000000)
	if err != nil {
		log.Error("unexpected error:", err)
		return
	}
	for i, w := range list {
		log.Infof("account[%d] : %+v", i, w)
	}
	log.Info("account count:", len(list))

	tm.CloseDB(testApp)

}

func TestWalletManager_CreateAddress(t *testing.T) {

	tm := testInitWalletManager()

	walletID := "W2DyYXbPCpkXWS1tJPYcRxhioSNyqwSu8F"
	//accountID := "9YBe43SkTyBneYNEnR7tHB3dh7VPB7toYkaZzU869C9y"
	accountID := "EhXYgY4wFN91VzkmJtyXPa1mPwEcp7o7PokQqaKcKGE4"
	address, err := tm.CreateAddress(testApp, walletID, accountID, 5)
	if err != nil {
		log.Error(err)
		return
	}

	for _, addr := range address {
		log.Info(addr.Address)
	}

	tm.CloseDB(testApp)
}

func TestWalletManager_GetAddressList(t *testing.T) {

	tm := testInitWalletManager()

	walletID := "W2DyYXbPCpkXWS1tJPYcRxhioSNyqwSu8F"
	//accountID := "9YBe43SkTyBneYNEnR7tHB3dh7VPB7toYkaZzU869C9y"
	accountID := "EhXYgY4wFN91VzkmJtyXPa1mPwEcp7o7PokQqaKcKGE4"
	//accountID := "47VD3c4xUuvCu1cuaQffRMcgQdkkAtYovUwwiMNFpKNe"
	list, err := tm.GetAddressList(testApp, walletID, accountID, 0, -1, false)
	if err != nil {
		log.Error("unexpected error:", err)
		return
	}
	for _, w := range list {
		fmt.Println(w.Address)
	}
	log.Info("address count:", len(list))

	tm.CloseDB(testApp)
}
