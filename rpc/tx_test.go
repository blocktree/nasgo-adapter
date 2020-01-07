package rpc

import (
	"reflect"
	"testing"
)

func TestTx_GetTransaction(t *testing.T) {
	type fields struct {
		baseAddress string
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Transaction
		wantErr bool
	}{
		{
			name:   "test get tx",
			fields: fields{baseAddress: Url},
			args: args{
				id: "da463370c7102425e3d4e9e0317d01efe192710f4614f4724abd8ceacbef8b1d",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.fields.baseAddress)
			got, err := client.Tx.GetTransaction(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tx.GetTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tx.GetTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTx_GetTransactions(t *testing.T) {
	type fields struct {
		baseAddress string
	}
	type args struct {
		blockId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Transaction
		wantErr bool
	}{
		{
			name:   "test get txs",
			fields: fields{baseAddress: Url},
			args: args{
				blockId: "ab2eb41872f70d8da38514110f6d24a12b7c4115401ba86be1b9a0c5c002c039",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.fields.baseAddress)
			got, err := client.Tx.GetTransactionsByBlock(tt.args.blockId)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tx.GetTransactionsByBlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tx.GetTransactionsByBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}
