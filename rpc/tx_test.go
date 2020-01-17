package rpc

import (
	"encoding/json"
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
				id: "e94bc52c095d7321fe7dca50c42ab3b4337bd06af3a80d958b0d7dcbade5c526",
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

func TestTx_BroadcastTx(t *testing.T) {
	rawTx := `
{"transaction":{"type":0,"amount":123456,"fee":1000000,"recipientId":"NDt9qnAHnFAuP8T9GbzQ2o8UaacQscAcU2","message":"hello boy","timestamp":59049090,"asset":{},"senderPublicKey":"1a43612ad299bc749395ac164878044d3aee89cedc8fed7a08f13e3ad1b4fedc","signature":"1bb98abf922aae37175a980fcda371c1dcacb9aea5d82ae0903f4886b3d5a8423da0641872ea58407b8b47ca7ce3d82b2538791e127c6fc9534ed161e8b2a008","id":"edda385d1824b9f28d5f78dcd54c8e8d6004182de48b2302d892fb5adc427c88"}}
`
	var tx map[string]interface{}
	err := json.Unmarshal([]byte(rawTx), &tx)
	if err != nil {
		t.Errorf("json.Unmarshal error = %v", err)
		return
	}

	client := NewClient("http://localhost:20001")
	err = client.Tx.BroadcastTx(tx, 1)
	if err != nil {
		t.Errorf("json.Unmarshal error = %v", err)
		return
	}
	//log.Infof("txid: %v", tx["id"])
}
