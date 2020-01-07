package rpc

import (
	"reflect"
	"testing"
)

func TestWallet_GetBalance(t *testing.T) {
	type fields struct {
		baseAddress string
	}
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    uint64
		wantErr bool
	}{
		{
			name:   "test balance",
			fields: fields{baseAddress: Url},
			args: args{
				address: "N6R2Gq5ZsXuMcXYcZo1bQ9WQRnaT9Ba3Ha",
			},
			want:    17200000,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.fields.baseAddress)
			got, err := client.Wallet.GetBalance(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("Wallet.GetBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Wallet.GetBalance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWallet_GetAssetsBalance(t *testing.T) {
	type fields struct {
		baseAddress string
	}
	type args struct {
		address  string
		currency string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *AssetsBalance
		wantErr bool
	}{
		{
			name:   "test assets balance",
			fields: fields{baseAddress: Url},
			args: args{
				address:  "N2BxCH4EJqD2ACr4vto6nFQML25oZ5LE2g",
				currency: "OBX.OBX",
			},
			want: &AssetsBalance{Currency: "OBX.OBX", Balance: "100000000",
				Precision: 8},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.fields.baseAddress)
			got, err := client.Wallet.GetAssetsBalance(tt.args.address, tt.args.currency)
			if (err != nil) != tt.wantErr {
				t.Errorf("Wallet.GetAssetsBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Wallet.GetAssetsBalance() = %v, want %v", got, tt.want)
			}
		})
	}
}
