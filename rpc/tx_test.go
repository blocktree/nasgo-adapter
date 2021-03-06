package rpc

import (
	"encoding/json"
	"github.com/blocktree/openwallet/v2/log"
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
				id: "51db69b4a4917dae6a230925777b4591ccd67c4bab367afcd9ef041c9183c72e",
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
				blockId: "b0938069b59f336482220a0128bf8b4874ed49792b354a2e74bafcd759a1bd15",
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
			for _, tx := range got {
				log.Infof("tx: %+v", tx)
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
