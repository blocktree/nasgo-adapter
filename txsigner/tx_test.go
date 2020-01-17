package txsigner

import (
	"encoding/hex"
	"github.com/blocktree/go-owcrypt"
	"github.com/blocktree/openwallet/log"
	"reflect"
	"testing"

	"github.com/blocktree/nasgo-adapter/rpc"
)

func TestTransaction_GenerateHash(t *testing.T) {
	pub := "d67925c8c7fda675b4bf8e3230d2fccafd9c32be6414059bc3aa4bbb87d88548"
	want, _ := hex.DecodeString("f549d5d0b6e67e0392205ed1e07bd8278333e82fd683988c4aed4d6e90e53bae")
	type fields struct {
		Transaction     *rpc.Transaction
		SenderPublicKey string
	}
	tests := []struct {
		name     string
		fields   fields
		wantHash []byte
	}{
		{
			name: "test hash",
			fields: fields{
				Transaction: &rpc.Transaction{
					Amount:      12345678,
					Fee:         1000000,
					Message:     "hello boy",
					RecipientId: "NDt9qnAHnFAuP8T9GbzQ2o8UaacQscAcU2",
					Timestamp:   58982624,
					Type:        0,
				},
				SenderPublicKey: pub,
			},
			wantHash: want,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := &Transaction{
				Transaction:     tt.fields.Transaction,
				SenderPublicKey: tt.fields.SenderPublicKey,
			}
			if gotHash := tx.GenerateHash(true); !reflect.DeepEqual(gotHash, tt.wantHash) {
				t.Errorf("Transaction.GenerateHash() = %v, want %v", hex.EncodeToString(gotHash), tt.wantHash)
			}
		})
	}
}

func TestVerify(t *testing.T) {
	hash, _ := hex.DecodeString("05fad362af2fc19203c8603b875f30397650edff91017a445261c59a0c18c416")
	pub, _ := hex.DecodeString("d67925c8c7fda675b4bf8e3230d2fccafd9c32be6414059bc3aa4bbb87d88548")
	signature, _ := hex.DecodeString("7d774086cf12dda36323aff2279a54d39ffd9c81ddf1549938170ea947b729133af84f92612f1fa6e2ad36eabf3d278a393b2483148db50a4ee2213cb8c22e07")
	flag := owcrypt.Verify(pub, nil, hash, signature, owcrypt.ECC_CURVE_ED25519)
	if flag == owcrypt.SUCCESS {
		log.Infof("success")
	} else {
		log.Infof("fail")
	}

}
