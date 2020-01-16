package txsigner

import (
	"encoding/hex"
	"fmt"

	"github.com/blocktree/go-owcrypt"
	"github.com/blocktree/nasgo-adapter/rpc"
	"github.com/blocktree/nasgo-adapter/utils"
)

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
	if tx.Type != rpc.TxType_NSG && tx.Type != rpc.TxType_Asset {
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
