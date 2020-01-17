package txsigner

import (
	"encoding/hex"
	"fmt"
	"github.com/blocktree/openwallet/log"

	"github.com/blocktree/go-owcrypt"
	"github.com/blocktree/nasgo-adapter/rpc"
	"github.com/blocktree/nasgo-adapter/utils"
)

type Transaction struct {
	*rpc.Transaction
	SenderPublicKey string `json:"senderPublicKey"`
}

func (tx *Transaction) GetID() string {
	hash := tx.GenerateHash(false)
	return hex.EncodeToString(hash)
}

func (tx *Transaction) GenerateHash(skipSignature bool) (hash []byte) {

	if tx == nil {
		fmt.Errorf("transaction is empty")
		return
	}
	if tx.Type != rpc.TxType_NSG && tx.Type != rpc.TxType_Asset {
		fmt.Errorf("transaction type is not allowed: %v", tx.Type)
		return
	}

	// signature, _ := hex.DecodeString(tx.Signature)

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

	pubBytes, _ := hex.DecodeString(tx.SenderPublicKey)

	txSlices := [][]byte{
		utils.PutUInt32ToBytes(tx.Type),
		utils.UInt32ToBytes(uint32(tx.Timestamp)),
		pubBytes,
		[]byte(tx.RecipientId),
		utils.UInt64ToBytes(tx.Amount),
		[]byte(tx.Message),
		assetSlice,
		// signature,
	}

	if !skipSignature && len(tx.Signature) > 0 {
		signature, _ := hex.DecodeString(tx.Signature)
		txSlices = append(txSlices, signature)
	}

	msg := utils.ConcatByteArray(txSlices)
	log.Debugf("tx msg: %s", hex.EncodeToString(msg))
	hash = owcrypt.Hash(msg, 0, owcrypt.HASH_ALG_SHA256)
	return
}
