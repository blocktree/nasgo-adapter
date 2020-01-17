package rpc

import (
	"reflect"
	"testing"
)

const Url = "http://localhost:20001"

func TestBlock_GetByHeight(t *testing.T) {
	type fields struct {
		baseAddress string
	}
	type args struct {
		height uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Header
		wantErr bool
	}{
		{
			name:   "Normal test",
			fields: fields{baseAddress: Url},
			args: args{
				height: 5093159,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.fields.baseAddress)
			got, err := client.Block.GetByHeight(tt.args.height)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Error("nil response")
			}
			if got != nil && reflect.DeepEqual(got.ID, [32]byte{}) {
				t.Error("empty hash")
			}
			t.Logf("%+v", *got)
		})
	}
}

func TestBlock_GetBlockHeight(t *testing.T) {
	type fields struct {
		baseAddress string
	}
	tests := []struct {
		name    string
		fields  fields
		want    uint32
		wantErr bool
	}{
		{
			name:    "Get Height",
			fields:  fields{baseAddress: Url},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.fields.baseAddress)
			got, err := client.Block.GetBlockHeight()
			if (err != nil) != tt.wantErr {
				t.Errorf("Block.GetBlockHeight() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == 0 {
				t.Errorf("Block.GetBlockHeight() = %v, want %v", got, tt.want)
			}
		})
	}
}
