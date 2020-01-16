package txsigner

import (
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/blocktree/nasgo-adapter/rpc"
)

func TestTransaction_GenerateHash(t *testing.T) {
	pub, _ := hex.DecodeString("d67925c8c7fda675b4bf8e3230d2fccafd9c32be6414059bc3aa4bbb87d88548")
	want, _ := hex.DecodeString("f549d5d0b6e67e0392205ed1e07bd8278333e82fd683988c4aed4d6e90e53bae")
	type fields struct {
		Transaction     *rpc.Transaction
		SenderPublicKey []byte
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
			if gotHash := tx.GenerateHash(); !reflect.DeepEqual(gotHash, tt.wantHash) {
				t.Errorf("Transaction.GenerateHash() = %v, want %v", gotHash, tt.wantHash)
			}
		})
	}
}
