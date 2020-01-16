package txsigner

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/blocktree/go-owcrypt"
)

var Default = &TransactionSigner{}

type TransactionSigner struct {
}

// SignTransactionHash 交易哈希签名算法
// required
func (singer *TransactionSigner) SignTransactionHash(msg []byte, privateKey []byte, eccType uint32) ([]byte, error) {
	signature, _, retCode := owcrypt.Signature(privateKey, nil, msg, owcrypt.ECC_CURVE_ED25519)
	if retCode != owcrypt.SUCCESS {
		return nil, fmt.Errorf("ECC sign hash failed")
	}

	return signature, nil
}

// VerifyAndCombineTransaction verify signature
// required
func (singer *TransactionSigner) VerifyAndCombineTransaction(emptyTrans string, message, signature, publickKey []byte) (bool, string, error) {
	trx := Transaction{}

	err := json.Unmarshal([]byte(emptyTrans), &trx)

	if err != nil {
		return false, "", errors.New("Invalid empty transaction data")
	}

	ret := owcrypt.Verify(publickKey, nil, message, signature, owcrypt.ECC_CURVE_ED25519)
	if ret != owcrypt.SUCCESS {
		errinfo := fmt.Sprintf("verify error, ret:%v\n", "0x"+strconv.FormatUint(uint64(ret), 16))
		return false, "", errors.New(errinfo)
	}

	trx.Signature = hex.EncodeToString(signature)
	txBytes, err := json.Marshal(trx)
	if err != nil {
		return false, "", errors.New("Failed to marshal transaction")
	}

	return true, hex.EncodeToString(txBytes), nil
}
