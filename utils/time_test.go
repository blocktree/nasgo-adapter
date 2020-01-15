package utils

import (
	"reflect"
	"testing"
	"time"
)

func TestGetEpochTime(t *testing.T) {
	tests := []struct {
		name string
		want int64
	}{
		{
			name: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEpochTime(); got != tt.want {
				t.Errorf("GetEpochTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTime(t *testing.T) {
	type args struct {
		t int64
	}
	ti, _ := time.Parse("2006-01-02 15:04 MST", "2020-01-14 13:56 UTC")
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "t",
			args: args{t: 58787815},
			want: ti,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTime(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
